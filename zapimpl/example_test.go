package zapimpl_test

import (
	"github.com/junk1tm/log"
	"go.uber.org/zap"

	"github.com/junk1tm/log/zapimpl"
)

func ExampleNewLogger() {
	// configure zap logger here:
	zl := zap.NewExample()

	logger := zapimpl.NewLogger(zl)
	logger.Debug("example 1", log.Int("foo", 1))
	logger.Info("example 2", log.Int("bar", 2))
	logger.Error("example 3", log.Int("baz", 3))

	// output:
	// {"level":"debug","msg":"example 1","foo":1}
	// {"level":"info","msg":"example 2","bar":2}
	// {"level":"error","msg":"example 3","baz":3}
}
