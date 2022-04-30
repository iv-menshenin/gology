package gology

import (
	"fmt"
	"io"
	"log"
	"sync/atomic"

	"github.com/pkg/errors"
)

const (
	poolSize          = 100
	defaultBufferSize = 1024
)

//nolint:gochecknoglobals
var pool = make(chan reusable, poolSize)

func newLogger() Logger {
	// allocations
	return Logger{
		reusable: reusable{
			buffer: make([]byte, 0, defaultBufferSize),
			closed: new(int64),
		},
	}
}

func acquireLogger() Logger {
	select {

	case r := <-pool:
		atomic.StoreInt64(r.closed, 0)
		return Logger{
			reusable: r,
		}

	default:
		return newLogger()
	}
}

func releaseLogger(r reusable) {
	atomic.StoreInt64(r.closed, 1)
	r.buffer = r.buffer[:0]

	select {

	case pool <- r:
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

func attrsToJSON(b []byte, attrs ...Attr) []byte {
	var needSep = len(b) > 0 && b[len(b)-1] != '{'
	for _, attr := range attrs {
		if needSep {
			b = append(b, ',')
		}
		needSep = true
		b = append(b, '"')
		b = append(b, attr.name...)
		b = append(b, '"', ':')
		b = attrToJSON(b, attr)
	}
	return b
}

func attrToJSON(b []byte, attr Attr) []byte {
	switch attr.typ {

	case attrString:
		b = strAttrToJSON(b, attr.str)

	case attrInt:
		b = intAttrToJSON(b, attr.int)

	case attrUint:
		b = uintAttrToJSON(b, attr.uint)

	case attrTime:
		// allocations
		b = strAttrToJSON(b, attr.tm.String())

	case attrError:
		// allocations
		type stackTracer interface {
			StackTrace() errors.StackTrace
		}
		if attr.err == nil {
			b = append(b, "null"...)
			break
		}

		b = strAttrToJSON(b, attr.err.Error())
		if t, ok := attr.err.(stackTracer); ok {
			b = append(b, ",\"stack\":"...)
			b = strAttrToJSON(b, fmt.Sprintf("%+v", t.StackTrace()))
		}

	}
	return b
}

func strAttrToJSON(b []byte, attr string) []byte {
	b = append(b, '"')
	b = safeStringAppend(b, attr)
	b = append(b, '"')
	return b
}

func intAttrToJSON(b []byte, attr int64) []byte {
	if attr == 0 {
		b = append(b, '0')
		return b
	}
	var start = attr
	var neg bool
	if attr < 0 {
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
	return b
}

func uintAttrToJSON(b []byte, attr uint64) []byte {
	if attr == 0 {
		b = append(b, '0')
		return b
	}
	var ib [22]byte
	var ip = 22
	for n := attr; n > 0; n = n / 10 {
		ip--
		ib[ip] = byte(rune('0' + n%10))
	}
	b = append(b, ib[ip:]...)
	return b
}

func safeStringAppend(b []byte, s string) []byte {
	if s == "" {
		return b
	}
	var l int
	for i, n := range s {
		switch n {

		case '"':
			b = append(b, s[l:i]...)
			b = append(b, '\\', '"')
			l = i + 1

		case '\n':
			b = append(b, s[l:i]...)
			b = append(b, '\\', 'n')
			l = i + 1

		case '\t':
			b = append(b, s[l:i]...)
			b = append(b, '\\', 't')
			l = i + 1
		}
	}
	return append(b, s[l:]...)
}
