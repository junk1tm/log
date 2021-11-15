// Package logrusimpl contains logrus implementation of Logger interface.
package logrusimpl

import (
	"github.com/junk1tm/log"
	"github.com/sirupsen/logrus"
)

// NewLogger creates a new log.Logger from the provided logrus.Logger.
func NewLogger(logger *logrus.Logger) log.Logger {
	return &wrapper{
		logger: logger,
	}
}

type wrapper struct {
	logger *logrus.Logger
}

func (w *wrapper) Debug(msg string, fields ...log.Field) {
	w.logger.WithFields(logrusFields(fields)).Debug(msg)
}

func (w *wrapper) Info(msg string, fields ...log.Field) {
	w.logger.WithFields(logrusFields(fields)).Info(msg)
}

func (w *wrapper) Error(msg string, fields ...log.Field) {
	w.logger.WithFields(logrusFields(fields)).Error(msg)
}

func logrusFields(fields []log.Field) map[string]interface{} {
	lf := make(map[string]interface{}, len(fields))

	for _, field := range log.FlattenFields(fields) {
		lf[field.Key] = field.Value
	}

	return lf
}

// Unwrap unwraps the provided logger,
// allowing access to the underlying logrus.Logger.
// It returns true on success, false otherwise.
func Unwrap(logger log.Logger) (*logrus.Logger, bool) {
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
