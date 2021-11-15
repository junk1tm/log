// Package stdlogimpl contains stdlib logger implementation of Logger interface.
package stdlogimpl

import (
	"fmt"
	stdlog "log"
	"strings"

	"github.com/junk1tm/log"
)

// NewLogger creates a new log.Logger from the provided stdlog.Logger.
// It prints the provided fields in "key=value" form after the message.
// NOTE: The prefix of a standard logger is used to print the logging level,
// so setting it before calling NewLogger will have no effect.
func NewLogger(logger *stdlog.Logger) log.Logger {
	return &wrapper{
		callerSkip: 1 + 1, // default calldepth + this wrapper
		logger:     logger,
	}
}

type wrapper struct {
	callerSkip int
	logger     *stdlog.Logger
}

func (w *wrapper) Debug(msg string, fields ...log.Field) { w.log("DEBUG", msg, fields) }
func (w *wrapper) Info(msg string, fields ...log.Field)  { w.log("INFO", msg, fields) }
func (w *wrapper) Error(msg string, fields ...log.Field) { w.log("ERROR", msg, fields) }

func (w *wrapper) AddCallerSkip(skip int) { w.callerSkip += skip }

func (w *wrapper) log(lvl string, msg string, fields []log.Field) {
	prefix := fmt.Sprintf("[%s] ", lvl)
	w.logger.SetPrefix(prefix)

	var sb strings.Builder
	sb.WriteString(msg)
	for _, field := range fields {
		_, _ = fmt.Fprintf(&sb, " %s=%v", field.Key, field.Value)
	}

	_ = w.logger.Output(w.callerSkip+1, sb.String())
}

// Unwrap unwraps the provided logger,
// allowing access to the underlying stdlog.Logger.
// It returns true on success, false otherwise.
func Unwrap(logger log.Logger) (*stdlog.Logger, bool) {
	for {
		switch l := logger.(type) {
		case *wrapper:
			return l.logger, true
		case interface{ Unwrap() log.Logger }:
			logger = l.Unwrap()
		default:
			return nil, false
		}
	}
}
