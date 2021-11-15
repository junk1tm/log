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
	var testLogger testLogger
	logger := log.WithFields(&testLogger, log.Int("foo", 1))
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
	if got := testLogger.calls; !reflect.DeepEqual(got, want) {
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

	var testLogger testLogger
	logger := log.WithHooks(&testLogger, prefixHook, multiplierHook, eofHook)

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
	if got := testLogger.calls; !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v; want %+v", got, want)
	}
}

type call struct {
	msg    string
	fields []log.Field
}

type testLogger struct {
	calls      []call
	callerSkip int
}

func (tl *testLogger) Debug(msg string, fields ...log.Field) {
	tl.calls = append(tl.calls, call{msg: msg, fields: append(fields, tl.callerField())})
}

func (tl *testLogger) Info(msg string, fields ...log.Field) {
	tl.calls = append(tl.calls, call{msg: msg, fields: append(fields, tl.callerField())})
}

func (tl *testLogger) Error(msg string, fields ...log.Field) {
	tl.calls = append(tl.calls, call{msg: msg, fields: append(fields, tl.callerField())})
}

func (tl *testLogger) AddCallerSkip(skip int) {
	tl.callerSkip += skip
}

func (tl *testLogger) callerField() log.Field {
	_, file, line, _ := runtime.Caller(tl.callerSkip + 2)
	file = file[strings.LastIndex(file, "/")+1:]
	value := fmt.Sprintf("%s:%d", file, line)

	return log.String("caller", value)
}
