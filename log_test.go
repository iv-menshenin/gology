package gology

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type nullWriter struct{}

func (w *nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func Benchmark_Logger(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var log = New(&nullWriter{}, LevelDebug)
		log.Write(
			LevelDebug,
			"some message to write into test message log hello",
			String("debug", "some string attr! LevelDebug"),
			Int("count", n+1024),
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
}

func Test_Attrs(t *testing.T) {

	t.Run("simple test", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		log.Error("something went wrong", Int("count", 2), Int("wants", 99), String("info", "foo bar"))
		var i map[string]interface{}
		if err := json.Unmarshal(writer.Bytes(), &i); err != nil {
			t.Error(err)
		}
		var want = map[string]interface{}{
			"level":   "ERROR",
			"message": "something went wrong",
			"count":   2,
			"wants":   99,
			"info":    "foo bar",
		}
		if !mapMatch(i, want) {
			t.Errorf("matching error\na1: %+v\na2: %+v", want, i)
		}
	})

	t.Run("integers", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		log.Error(
			"something went wrong",
			Int("i", 2),
			Int16("i16", -333),
			Int32("i32", 43564),
			Int64("i64", 654456),
		)
		var i map[string]interface{}
		if err := json.Unmarshal(writer.Bytes(), &i); err != nil {
			t.Error(err)
		}
		var want = map[string]interface{}{
			"level":   "ERROR",
			"message": "something went wrong",
			"i":       2,
			"i16":     -333,
			"i32":     43564,
			"i64":     654456,
		}
		if !mapMatch(i, want) {
			t.Errorf("matching error\na1: %+v\na2: %+v", want, i)
		}
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
		var i map[string]interface{}
		if err := json.Unmarshal(writer.Bytes(), &i); err != nil {
			t.Error(err)
		}
		var want = map[string]interface{}{
			"level":   "ERROR",
			"message": "something went wrong",
			"i":       2,
			"i16":     333,
			"i32":     43564,
			"i64":     654456,
		}
		if !mapMatch(i, want) {
			t.Errorf("matching error\na1: %+v\na2: %+v", want, i)
		}
	})

	t.Run("date-time", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		var tm = time.Now().UTC().Round(time.Second)
		log.Error(
			"something went wrong",
			DateTime("date", tm),
		)
		var i map[string]interface{}
		if err := json.Unmarshal(writer.Bytes(), &i); err != nil {
			t.Error(err)
		}
		var want = map[string]interface{}{
			"level":   "ERROR",
			"message": "something went wrong",
			"date":    tm.String(),
		}
		if !mapMatch(i, want) {
			t.Errorf("matching error\na1: %+v\na2: %+v", want, i)
		}
	})

}

func Test_WithCascade(t *testing.T) {
	t.Run("with", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		var tm = time.Now().UTC().Round(time.Second)
		log1 := log.WithAttrs(DateTime("started", tm), String("test", "foo"))
		log1.Debug("test message", Int("int", 445))

		var i map[string]interface{}
		if err := json.Unmarshal(writer.Bytes(), &i); err != nil {
			t.Error(err)
		}
		var want = map[string]interface{}{
			"level":   "DEBUG",
			"test":    "foo",
			"message": "test message",
			"started": tm.String(),
			"int":     445,
		}
		if !mapMatch(i, want) {
			t.Errorf("matching error\na1: %+v\na2: %+v", want, i)
		}
	})

	t.Run("with cascade", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		var tm = time.Now().UTC().Round(time.Second)

		log1 := log.WithAttrs(DateTime("started", tm), String("test1", "foo"))
		log2 := log1.WithAttrs(String("test2", "bar"))
		log2.Debug("test message", Int("int", 445))

		var i map[string]interface{}
		if err := json.Unmarshal(writer.Bytes(), &i); err != nil {
			t.Error(err)
		}
		var want = map[string]interface{}{
			"level":   "DEBUG",
			"test1":   "foo",
			"test2":   "bar",
			"message": "test message",
			"started": tm.String(),
			"int":     445,
		}
		if !mapMatch(i, want) {
			t.Errorf("matching error\na1: %+v\na2: %+v", want, i)
		}
	})

	t.Run("with cancelled cascade", func(t *testing.T) {
		var writer = bytes.NewBuffer(make([]byte, 0, 65536))
		var log = New(writer, LevelDebug)
		var tm = time.Now().UTC().Round(time.Second)

		log1 := log.WithAttrs(DateTime("started", tm), String("test1", "foo"))
		_ = log1.WithAttrs(String("test2", "bar"))
		log1.Debug("test message", Int("int", 445))

		var i map[string]interface{}
		if err := json.Unmarshal(writer.Bytes(), &i); err != nil {
			t.Error(err)
		}
		var want = map[string]interface{}{
			"level":   "DEBUG",
			"test1":   "foo",
			"message": "test message",
			"started": tm.String(),
			"int":     445,
		}
		if !mapMatch(i, want) {
			t.Errorf("matching error\na1: %+v\na2: %+v", want, i)
		}
	})
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
