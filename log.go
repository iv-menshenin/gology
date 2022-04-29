package gology

import (
	"io"
	"log"
)

var bufferPool = make(chan *buffer, 10)

type buffer struct {
	b []byte
}

func newBuff() *buffer {
	return &buffer{
		b: make([]byte, 0, 1024),
	}
}

func acquireBuff() *buffer {
	select {

	case a := <-bufferPool:
		return a

	default:
		return newBuff()
	}
}

func releaseBuff(b *buffer) {
	b.b = b.b[:0]

	select {

	case bufferPool <- b:
		return

	default:
		// lost buffer
	}
}

func New(writer io.Writer, level Level) *Logger {
	return &Logger{
		writer:   writer,
		buffer:   acquireBuff(),
		logLevel: level,
	}
}

type (
	Logger struct {
		writer   io.Writer
		buffer   *buffer
		logLevel Level
	}
	Attr struct {
		name string
		int  int64
		str  string
	}
	Level int8
)

const (
	LevelError Level = iota
	LevelWarning
	LevelDebug
)

func (l *Logger) Write(level Level, message string, attrs ...Attr) {
	if level > l.logLevel {
		return
	}
	l.buffer.b = append(l.buffer.b, "{\"message\":\""...)
	l.buffer.b = append(l.buffer.b, message...)
	l.buffer.b = append(l.buffer.b, "\",\"level\":\""...)
	l.buffer.b = levelToBytes(l.buffer.b, level)
	l.buffer.b = append(l.buffer.b, '"')
	l.buffer.b = attrsToJson(l.buffer.b, attrs...)
	l.buffer.b = append(l.buffer.b, '}')
	if _, err := l.writer.Write(l.buffer.b); err != nil {
		log.Printf("cannot write message to logfile: %v", err)
	}
	l.buffer.b = l.buffer.b[:0]
}

func levelToBytes(b []byte, l Level) []byte {
	switch l {
	case LevelError:
		return append(b, "ERROR"...)
	case LevelWarning:
		return append(b, "WARNING"...)
	case LevelDebug:
		return append(b, "DEBUG"...)
	}
	return append(b, "UNKNOWN"...)
}

func (l *Logger) Error(message string, attrs ...Attr) {
	l.Write(LevelError, message, attrs...)
}

func (l *Logger) Warning(message string, attrs ...Attr) {
	l.Write(LevelWarning, message, attrs...)
}

func (l *Logger) Debug(message string, attrs ...Attr) {
	l.Write(LevelDebug, message, attrs...)
}

func Int(name string, i int) Attr {
	return Attr{name: name, int: int64(i)}
}

func String(name, s string) Attr {
	return Attr{name: name, str: s}
}

func attrsToJson(b []byte, attrs ...Attr) []byte {
	for _, attr := range attrs {
		b = append(b, ',', '"')
		b = append(b, attr.name...)
		b = append(b, '"', ':')

		switch {

		case attr.str != "":
			b = append(b, '"')
			b = append(b, attr.str...)
			b = append(b, '"')

		case attr.int != 0:
			var ib [20]byte
			var ip = 20
			for n := attr.int; n > 0; n = n / 10 {
				ip--
				ib[ip] = byte(rune('0' + n%10))
			}
			b = append(b, ib[ip:]...)
		}
	}
	return b
}

func (l *Logger) Close() {
	releaseBuff(l.buffer)
}
