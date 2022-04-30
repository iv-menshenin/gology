package gology

import (
	"testing"
)

type nullWriter struct{}

func (w *nullWriter) Write(p []byte) (n int, err error) {
	// fmt.Println(string(p))
	return len(p), nil
}

// CLEAR
// Benchmark_Logger-8   	 6390549	       187.8 ns/op	      32 B/op	       1 allocs/op

// WriteOut and JSON
// Benchmark_Logger-8   	 3752178	       310.3 ns/op	      32 B/op	       1 allocs/op

// NESTED
// Benchmark_Logger-8   	 3735806	       330.1 ns/op	      96 B/op	       3 allocs/op

func Benchmark_Logger(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var log = New(&nullWriter{}, LevelDebug)
		log.Write(
			LevelDebug,
			"some message to write into test message log hello",
			String("debug", "some string attr! LevelDebug"),
			Int("count", n+1024),
		)
		func(log *Logger) {
			log.Write(
				LevelError,
				"some message to write into test message log hello",
				String("error", "some string attr! LevelError"),
				Int("count", n+1024),
			)
			log.Write(
				LevelWarning,
				"repeat warning in nested",
				Int("test", 19),
			)
		}(log.WithAttrs(
			String("test1", "nested value 1"),
			String("test2", "nested value 2"),
		))
		log.Write(
			LevelWarning,
			"some message to write into test message log hello",
			String("warning", "some string attr! LevelWarning"),
			Int("count", n+1024),
		)
		log.Close()
	}
}
