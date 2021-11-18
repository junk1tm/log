package log_test

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/junk1tm/log"
)

func TestWithFields(t *testing.T) {
	var spy spyLogger
	logger := log.WithFields(&spy, log.Int("foo", 1))
	logger.Info("first call", log.Int("bar", 2))
	logger.Info("second call", log.Int("baz", 3))

	want := []call{
		{
			msg:    "first call",
			fields: []log.Field{log.Int("foo", 1), log.Int("bar", 2), log.String("caller", "log_test.go:18")},
		},
		{
			msg:    "second call",
			fields: []log.Field{log.Int("foo", 1), log.Int("baz", 3), log.String("caller", "log_test.go:19")},
		},
	}
	if got := spy.calls; !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v; want %+v", got, want)
	}
}

func TestWithHooks(t *testing.T) {
	// adds _ prefix to each key.
	prefixHook := func(lvl log.Level, msg string, fields []log.Field) error {
		for i := range fields {
			fields[i].Key = "_" + fields[i].Key
		}
		return nil
	}

	// multiplies each int value by 2.
	multiplierHook := func(lvl log.Level, msg string, fields []log.Field) error {
		for i := range fields {
			if v, ok := fields[i].Value.(int); ok {
				fields[i].Value = v * 2
			}
		}
		return nil
	}

	// fails with io.EOF but only at DEBUG level.
	eofHook := func(lvl log.Level, msg string, fields []log.Field) error {
		if lvl == log.DebugLevel {
			return io.EOF
		}
		return nil
	}

	log.OnHookError = func(logger log.Logger, err error) {
		if !errors.Is(err, io.EOF) {
			t.Errorf("got %v; want %v", err, io.EOF)
		}
	}

	var spy spyLogger
	logger := log.WithHooks(&spy, prefixHook, multiplierHook, eofHook)

	logger.Debug("first call", log.Int("foo", 1))
	logger.Info("second call", log.Int("bar", 2))

	want := []call{
		{
			msg:    "first call",
			fields: []log.Field{log.Int("_foo", 2), log.String("caller", "log_test.go:72")},
		},
		{
			msg:    "second call",
			fields: []log.Field{log.Int("_bar", 4), log.String("caller", "log_test.go:73")},
		},
	}
	if got := spy.calls; !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v; want %+v", got, want)
	}
}

// https://github.com/junk1tm/log/issues/11
func TestIssue11(t *testing.T) {
	// adds _ prefix to each key.
	prefixHook := func(lvl log.Level, msg string, fields []log.Field) error {
		for i := range fields {
			fields[i].Key = "_" + fields[i].Key
		}
		return nil
	}

	var spy spyLogger
	logger := log.WithHooks(&spy, prefixHook)
	logger = log.WithFields(logger, log.String("key", "value"))

	logger.Debug("first call")
	logger.Info("second call")
	logger.Error("third call")

	want := []call{
		{
			msg:    "first call",
			fields: []log.Field{log.String("_key", "value"), log.String("caller", "log_test.go:104")},
		},
		{
			msg:    "second call",
			fields: []log.Field{log.String("_key", "value"), log.String("caller", "log_test.go:105")},
		},
		{
			msg:    "third call",
			fields: []log.Field{log.String("_key", "value"), log.String("caller", "log_test.go:106")},
		},
	}
	if got := spy.calls; !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v; want %+v", got, want)
	}
}

type call struct {
	msg    string
	fields []log.Field
}

// spyLogger records its calls for later inspection in tests.
type spyLogger struct {
	calls      []call
	callerSkip int
}

func (sl *spyLogger) Debug(msg string, fields ...log.Field) {
	sl.calls = append(sl.calls, call{msg: msg, fields: append(fields, sl.callerField())})
}

func (sl *spyLogger) Info(msg string, fields ...log.Field) {
	sl.calls = append(sl.calls, call{msg: msg, fields: append(fields, sl.callerField())})
}

func (sl *spyLogger) Error(msg string, fields ...log.Field) {
	sl.calls = append(sl.calls, call{msg: msg, fields: append(fields, sl.callerField())})
}

func (sl *spyLogger) AddCallerSkip(skip int) {
	sl.callerSkip += skip
}

func (sl *spyLogger) callerField() log.Field {
	_, file, line, _ := runtime.Caller(sl.callerSkip + 2)
	file = file[strings.LastIndex(file, "/")+1:]
	value := fmt.Sprintf("%s:%d", file, line)

	return log.String("caller", value)
}
