package gology

import (
	"testing"
)

type (
	nullWriter struct{}
)

func (nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// CLEAR
// Benchmark_Logger-8   	 6390549	       187.8 ns/op	      32 B/op	       1 allocs/op

// WriteOut and JSON
// Benchmark_Logger-8   	 3752178	       310.3 ns/op	      32 B/op	       1 allocs/op

func Benchmark_Logger(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var log = New(nullWriter{}, LevelDebug)
		log.Write(
			LevelDebug,
			"some message to write into test message log hello",
			String("message", "some string attr! LevelDebug"),
			Int("count", n+1024),
		)
		log.Write(
			LevelWarning,
			"some message to write into test message log hello",
			String("message", "some string attr! LevelWarning"),
			Int("count", n+1024),
		)
		log.Write(
			LevelError,
			"some message to write into test message log hello",
			String("message", "some string attr! LevelError"),
			Int("count", n+1024),
		)
		log.Close()
	}
}
