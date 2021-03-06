package gology

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"

	errors "github.com/pkg/errors"
	// "github.com/rs/zerolog"
)

type nullWriter struct{}

func (w *nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

/*
	cpu: Intel(R) Core(TM) i7-9700F CPU @ 3.00GHz
	Benchmark_GologyLogger-8         6980311               171.7 ns/op             0 B/op          0 allocs/op
	Benchmark_ZeroLogger-8           1561195               876.7 ns/op           512 B/op          1 allocs/op
*/

func Benchmark_GologyLogger(b *testing.B) {
	var log = New(&nullWriter{}, LevelAll)
	for n := 0; n < b.N; n++ {
		log.Write(
			LevelDebug,
			"some message to write into test message log hello",
			String("debug", "some string attr! LevelDebug"),
			Int("count", n+1024),
			Float32("rate", 56.66),
		)
		func(log Logger) {
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
				Float64("rate", math.MaxFloat64),
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

//func Benchmark_ZeroLogger(b *testing.B) {
//	var log = zerolog.New(&nullWriter{})
//
//	for n := 0; n < b.N; n++ {
//		log.Debug().
//			Str("debug", "some string attr! LevelDebug").
//			Int("count", n+1024).
//			Msg("some message to write into test message log hello")
//
//		func(log zerolog.Logger) {
//			log.Error().
//				Str("debug", "some string attr! LevelDebug").
//				Int("count", n+1024).
//				Msg("some message to write into test message log hello")
//
//			log.Warn().
//				Str("debug", "some string attr! LevelDebug").
//				Int("test", 19).
//				Msg("repeat warning in nested")
//		}(
//			log.With().
//				Str("test1", "nested value 1").
//				Str("test2", "nested value 2").
//				Logger(),
//		)
//
//		log.Warn().
//			Str("warning", "some string attr! LevelWarning").
//			Int("count", n+1024).
//			Msg("some message to write into test message log hello")
//
//	}
//}

func Test_Levels(t *testing.T) {
	var writer = bytes.NewBuffer(make([]byte, 0, 65536))
	var log = New(writer, LevelDebug)
	var i map[string]interface{}

	t.Run("error", func(t *testing.T) {
		log.Error("something went wrong", Int("count", 2), Int("wants", 99), String("info", "foo bar"))
		if err := json.Unmarshal(writer.Bytes(), &i); err != nil {
			t.Error(err)
		}
		if lv := i["level"]; lv != "ERROR" {
			t.Errorf("level matching error\nneed: %+v\n got: %+v", "ERROR", lv)
		}
		writer.Reset()
	})

	t.Run("warning", func(t *testing.T) {
		log.Warning("something went wrong", Int("count", 2), Int("wants", 99), String("info", "foo bar"))
		if err := json.Unmarshal(writer.Bytes(), &i); err != nil {
			t.Error(err)
		}
		if lv := i["level"]; lv != "WARNING" {
			t.Errorf("level matching error\nneed: %+v\n got: %+v", "WARNING", lv)
		}
		writer.Reset()
	})

	t.Run("debug", func(t *testing.T) {
		log.Debug("something went wrong", Int("count", 2), Int("wants", 99), String("info", "foo bar"))
		if err := json.Unmarshal(writer.Bytes(), &i); err != nil {
			t.Error(err)
		}
		if lv := i["level"]; lv != "DEBUG" {
			t.Errorf("level matching error\nneed: %+v\n got: %+v", "DEBUG", lv)
		}
		writer.Reset()
	})

	t.Run("low level", func(t *testing.T) {
		var w = bytes.NewBuffer(make([]byte, 0, 65536))
		var l = New(w, LevelError)
		l.Warning("something warned")
		if w.Len() > 0 {
			t.Errorf("unexpected out: %s", w.String())
		}
	})

	log.Close()
}

func Test_Attrs(t *testing.T) {

	t.Run("zero test", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		log.Error(
			"something went wrong",
			UInt("count", 0),
			Int("max", 100),
			Int("wants", 0),
			String("info", ""),
			Float64("float", 0),
			Err(nil),
		)
		testResult(t, writer.Bytes(), map[string]interface{}{
			"level":   "ERROR",
			"message": "something went wrong",
			"count":   0,
			"max":     100,
			"wants":   0,
			"info":    "",
			"float":   0,
			"error":   nil,
		})
		log.Close()
	})

	t.Run("error test", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		log.Error(
			"something went wrong",
			Err(errors.WithStack(errors.New("error"))),
		)
		var i map[string]interface{}
		var got = writer.Bytes()
		if err := json.Unmarshal(got, &i); err != nil {
			t.Fatalf("error unmarshalling: %s\ndata: %s", err, string(got))
		}
		if stack, ok := i["stack"]; !ok && stack.(string) == "" {
			t.Fatalf("expected stack, got: %v", i["stack"])
		}
		log.Close()
	})

	t.Run("integers", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		log.Error(
			"something went wrong",
			Int("i", 2),
			Int16("i16", -333),
			Int32("i32", 43500),
			Int64("i64", 654456),
		)
		testResult(t, writer.Bytes(), map[string]interface{}{
			"level":   "ERROR",
			"message": "something went wrong",
			"i":       2,
			"i16":     -333,
			"i32":     43500,
			"i64":     654456,
		})
		log.Close()
	})

	t.Run("unsigned", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		log.Error(
			"something went wrong",
			UInt("i", 2),
			UInt16("i16", 333),
			UInt32("i32", 43564),
			UInt64("i64", 654456),
		)
		testResult(t, writer.Bytes(), map[string]interface{}{
			"level":   "ERROR",
			"message": "something went wrong",
			"i":       2,
			"i16":     333,
			"i32":     43564,
			"i64":     654456,
		})
		log.Close()
	})

	t.Run("error", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		log.Error(
			"something went wrong",
			Err(fmt.Errorf("fail io operations")),
		)
		testResult(t, writer.Bytes(), map[string]interface{}{
			"level":   "ERROR",
			"message": "something went wrong",
			"error":   "fail io operations",
		})
		log.Close()
	})

	t.Run("date-time", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		var tm = time.Now().UTC().Round(time.Second)
		log.Error(
			"something went wrong",
			DateTime("date", tm),
		)
		testResult(t, writer.Bytes(), map[string]interface{}{
			"level":   "ERROR",
			"message": "something went wrong",
			"date":    tm.String(),
		})
		log.Close()
	})

	t.Run("floats", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		log.Error(
			"something went wrong",
			Float64("f1", 999.999),
			Float32("f2", 999.999),
			Float64("f3", 1.123456789),
			Float64("f4", math.MaxFloat64),
		)
		testResult(t, writer.Bytes(), map[string]interface{}{
			"level":   "ERROR",
			"message": "something went wrong",
			"f1":      999.999,
			"f2":      999.999,
			"f3":      1.1234,
			"f4":      math.MaxFloat64,
		})
		log.Close()
	})

	t.Run("bad strings", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		log.Error("something \"went\" wrong", String("info", "foo \"bar\""))
		testResult(t, writer.Bytes(), map[string]interface{}{
			"level":   "ERROR",
			"message": "something \"went\" wrong",
			"info":    "foo \"bar\"",
		})
		log.Close()
	})

}

func Test_WithCascade(t *testing.T) {
	t.Run("with", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		var tm = time.Now().UTC().Round(time.Second)
		log1 := log.WithAttrs(DateTime("started", tm), String("test", "foo"))
		log1.Debug("test message", Int("int", 445))

		testResult(t, writer.Bytes(), map[string]interface{}{
			"level":   "DEBUG",
			"test":    "foo",
			"message": "test message",
			"started": tm.String(),
			"int":     445,
		})
		log.Close()
	})

	t.Run("with cascade", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		var tm = time.Now().UTC().Round(time.Second)

		log1 := log.WithAttrs(DateTime("started", tm), String("test1", "foo"))
		log2 := log1.WithAttrs(String("test2", "bar"))
		log2.Debug("test message", Int("int", 445))

		testResult(t, writer.Bytes(), map[string]interface{}{
			"level":   "DEBUG",
			"test1":   "foo",
			"test2":   "bar",
			"message": "test message",
			"started": tm.String(),
			"int":     445,
		})
		log.Close()
	})

	t.Run("with cancelled cascade", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		var tm = time.Now().UTC().Round(time.Second)

		log1 := log.WithAttrs(DateTime("started", tm), String("test1", "foo"))
		_ = log1.WithAttrs(String("test2", "bar"))
		log1.Debug("test message", Int("int", 445))

		testResult(t, writer.Bytes(), map[string]interface{}{
			"level":   "DEBUG",
			"test1":   "foo",
			"message": "test message",
			"started": tm.String(),
			"int":     445,
		})
		log.Close()
	})
}

func Test_SequencedWrite(t *testing.T) {
	var writer = bytes.NewBuffer(make([]byte, 0, 65536))
	var log = New(writer, LevelDebug).WithAttrs(String("foo", "bar"))
	var i map[string]interface{}

	t.Run("first out", func(t *testing.T) {
		log.Error("something went test", Int("count", 2))
		if err := json.Unmarshal(writer.Bytes(), &i); err != nil {
			t.Error(err)
		}
		if lv := i["level"]; lv != "ERROR" {
			t.Errorf("level matching error\nneed: %+v\n got: %+v", "ERROR", lv)
		}
		testResult(t, writer.Bytes(), map[string]interface{}{
			"message": "something went test",
			"level":   "ERROR",
			"foo":     "bar",
			"count":   2,
		})
		writer.Reset()
	})

	t.Run("second out", func(t *testing.T) {
		log.Error("something went test 2", Int("iter_cnt", 12))
		if err := json.Unmarshal(writer.Bytes(), &i); err != nil {
			t.Error(err)
		}
		if lv := i["level"]; lv != "ERROR" {
			t.Errorf("level matching error\nneed: %+v\n got: %+v", "ERROR", lv)
		}
		testResult(t, writer.Bytes(), map[string]interface{}{
			"message":  "something went test 2",
			"level":    "ERROR",
			"foo":      "bar",
			"iter_cnt": 12,
		})
		writer.Reset()
	})

	var log1 = log.WithAttrs(String("second", "foo bar bar foo"))

	t.Run("second out", func(t *testing.T) {
		log1.Debug("something went test 3", Int("free", 99))
		if err := json.Unmarshal(writer.Bytes(), &i); err != nil {
			t.Error(err)
		}
		testResult(t, writer.Bytes(), map[string]interface{}{
			"message": "something went test 3",
			"second":  "foo bar bar foo",
			"level":   "DEBUG",
			"foo":     "bar",
			"free":    99,
		})
		writer.Reset()
	})

	t.Run("quit out", func(t *testing.T) {
		log.Debug("something went test 4", Int("stock", 1))
		if err := json.Unmarshal(writer.Bytes(), &i); err != nil {
			t.Error(err)
		}
		testResult(t, writer.Bytes(), map[string]interface{}{
			"message": "something went test 4",
			"level":   "DEBUG",
			"foo":     "bar",
			"stock":   1,
		})
		writer.Reset()
	})
}

func Test_CloseLog(t *testing.T) {
	var writer = bytes.NewBuffer(make([]byte, 0, 65536))
	var log = New(writer, LevelDebug).WithAttrs(String("foo", "bar"))
	log.Close()
	log.Error("oops")
	if writer.Len() > 0 {
		t.Errorf("unexpected out: %s", writer.String())
	}

	log = New(writer, LevelDebug).WithAttrs(String("foo", "bar"))
	log.Error("oops")
	if writer.Len() == 0 {
		t.Error("unexpected empty out")
	}
}

func testResult(t *testing.T, got []byte, want map[string]interface{}) {
	t.Helper()
	var i map[string]interface{}
	if err := json.Unmarshal(got, &i); err != nil {
		t.Fatalf("error unmarshalling: %s\ndata: %s", err, string(got))
	}
	if !mapMatch(i, want) {
		t.Errorf("matching error\n extra got: %+v\nextra want: %+v", i, want)
	}
}

func mapMatch(a, b map[string]interface{}) bool {
	var remKeys []string
	for k, v1 := range a {
		if v2, ok := b[k]; ok && matchElement(v1, v2) {
			remKeys = append(remKeys, k)
		}
	}
	for _, key := range remKeys {
		delete(a, key)
		delete(b, key)
	}
	return len(a) == 0 && len(b) == 0
}

func matchElement(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
