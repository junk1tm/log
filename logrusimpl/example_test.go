package logrusimpl_test

import (
	"os"

	"github.com/junk1tm/log"
	"github.com/sirupsen/logrus"

	"github.com/junk1tm/log/logrusimpl"
)

func ExampleNewLogger() {
	// configure logrus logger here:
	ll := logrus.New()
	ll.Out = os.Stdout
	ll.Level = logrus.DebugLevel
	ll.Formatter = &logrus.JSONFormatter{DisableTimestamp: true}

	logger := logrusimpl.NewLogger(ll)
	logger.Debug("example 1", log.Int("foo", 1))
	logger.Info("example 2", log.Int("bar", 2))
	logger.Error("example 3", log.Int("baz", 3))

	// output:
	// {"foo":1,"level":"debug","msg":"example 1"}
	// {"bar":2,"level":"info","msg":"example 2"}
	// {"baz":3,"level":"error","msg":"example 3"}
}
