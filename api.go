package gology

import (
	"io"
	"sync/atomic"
	"time"
)

type (
	Logger struct {
		reusable
		writer   io.Writer
		logLevel Level
		fixed    int
	}
	reusable struct {
		buffer []byte
		closed *int64
	}
	Attr struct {
		name string
		int  int64
		uint uint64
		str  string
		tm   time.Time
	}
	Level int8
)

const (
	LevelError Level = iota
	LevelWarning
	LevelDebug
)

func New(writer io.Writer, level Level) Logger {
	var l = acquireLogger()
	l.writer = writer
	l.logLevel = level
	l.buffer = append(l.buffer, '{')
	return l
}

func (l Logger) Write(level Level, message string, attrs ...Attr) {
	if atomic.LoadInt64(l.closed) != 0 {
		return
	}
	if level > l.logLevel {
		return
	}
	var writer = l.writer
	var buffer = l.buffer
	if len(buffer) > 0 && buffer[len(buffer)-1] != '{' {
		buffer = append(buffer, ',')
	}
	buffer = append(buffer, "\"message\":\""...)
	buffer = safeStringAppend(buffer, message)
	buffer = append(buffer, "\",\"level\":\""...)
	buffer = levelToBytes(buffer, level)
	buffer = append(buffer, '"')
	buffer = attrsToJson(buffer, attrs...)
	buffer = append(buffer, '}')
	write(writer, buffer)
}

func (l Logger) Error(message string, attrs ...Attr) {
	l.Write(LevelError, message, attrs...)
}

func (l Logger) Warning(message string, attrs ...Attr) {
	l.Write(LevelWarning, message, attrs...)
}

func (l Logger) Debug(message string, attrs ...Attr) {
	l.Write(LevelDebug, message, attrs...)
}

func (l Logger) WithAttrs(attrs ...Attr) Logger {
	l.buffer = attrsToJson(l.buffer, attrs...)
	return l
}

func (l Logger) Close() {
	releaseLogger(l.reusable)
}
