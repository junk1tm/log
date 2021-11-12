package zerologimpl_test

import (
	"os"

	"github.com/junk1tm/log"
	"github.com/rs/zerolog"

	"github.com/junk1tm/log/zerologimpl"
)

func ExampleNewLogger() {
	// configure zerolog logger here:
	zl := zerolog.New(os.Stdout)

	logger := zerologimpl.NewLogger(zl)
	logger.Debug("example 1", log.Int("foo", 1))
	logger.Info("example 2", log.Int("bar", 2))
	logger.Error("example 3", log.Int("baz", 3))

	// output:
	// {"level":"debug","foo":1,"message":"example 1"}
	// {"level":"info","bar":2,"message":"example 2"}
	// {"level":"error","baz":3,"message":"example 3"}
}
