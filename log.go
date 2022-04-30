package gology

import (
	"io"
	"log"
)

const (
	poolSize          = 100
	defaultBufferSize = 1024
)

var pool = make(chan []byte, poolSize)

func newLogger() Logger {
	return Logger{
		buffer: make([]byte, 0, defaultBufferSize), // alloc
	}
}

func acquireLogger() Logger {
	select {

	case buffer := <-pool:
		return Logger{
			buffer: buffer,
		}

	default:
		return newLogger()
	}
}

func releaseLogger(buffer []byte) {

	select {

	case pool <- buffer:
		return

	default:
		// lost buffer
	}
}

func write(writer io.Writer, message []byte) {
	if _, err := writer.Write(message); err != nil {
		log.Printf("cannot write message to logfile: %v", err)
	}
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

func attrsToJson(b []byte, attrs ...Attr) []byte {
	var needSep = len(b) > 0 && b[len(b)-1] != '{'
	for _, attr := range attrs {
		if needSep {
			b = append(b, ',')
		}
		needSep = true
		b = append(b, '"')
		b = append(b, attr.name...)
		b = append(b, '"', ':')

		switch {

		case attr.str != "":
			b = append(b, '"')
			b = append(b, attr.str...)
			b = append(b, '"')

		case attr.int != 0:
			var start = attr.int
			var neg bool
			if attr.int < 0 {
				neg = true
				start = start * -1
			}
			var ib [22]byte
			var ip = 22
			for n := start; n > 0; n = n / 10 {
				ip--
				ib[ip] = byte(rune('0' + n%10))
			}
			if neg {
				b = append(b, '-')
			}
			b = append(b, ib[ip:]...)

		case attr.uint != 0:
			var ib [22]byte
			var ip = 22
			for n := attr.uint; n > 0; n = n / 10 {
				ip--
				ib[ip] = byte(rune('0' + n%10))
			}
			b = append(b, ib[ip:]...)

		case !attr.tm.IsZero():
			b = append(b, '"')
			b = append(b, attr.tm.String()...)
			b = append(b, '"')

		}

	}
	return b
}
