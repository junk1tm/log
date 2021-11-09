package log_test

import (
	"fmt"
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
			fields: []log.Field{log.Int("foo", 1), log.Int("bar", 2), log.String("caller", "log_test.go:16")},
		},
		{
			msg:    "second call",
			fields: []log.Field{log.Int("foo", 1), log.Int("baz", 3), log.String("caller", "log_test.go:17")},
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
