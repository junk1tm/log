package stdlogimpl_test

import (
	stdlog "log"
	"os"

	"github.com/junk1tm/log"

	"github.com/junk1tm/log/stdlogimpl"
)

func ExampleNewLogger() {
	// configure stdlog logger here:
	sl := stdlog.New(os.Stdout, "", 0)

	logger := stdlogimpl.NewLogger(sl)
	logger.Debug("example 1", log.Int("foo", 1))
	logger.Info("example 2", log.Int("bar", 2))
	logger.Error("example 3", log.Int("baz", 3))

	// output:
	// [DEBUG] example 1 foo=1
	// [INFO] example 2 bar=2
	// [ERROR] example 3 baz=3
}

func ExampleUnwrap() {
	sl := stdlog.New(os.Stdout, "", 0)
	logger := stdlogimpl.NewLogger(sl)
	if _, ok := stdlogimpl.Unwrap(logger); ok {
		// use stdlog logger here:
	}
}
