// Package log provides a general-purpose logging API.
package log

import (
	"time"
)

// Logger is a minimal logging interface with levels.
// It uses structured logging with Field parameters to ensure type safety.
type Logger interface {
	// Debug writes logs at DEBUG level.
	// It is used to log information useful for developers.
	Debug(msg string, fields ...Field)
	// Info writes logs at INFO level.
	// It is used to log information useful for users.
	Info(msg string, fields ...Field)
	// Error writes logs at ERROR level.
	// It is used to handle errors by logging them.
	Error(msg string, fields ...Field)
}

// Loggable represents an entity that can be logged.
// It provides a type-safe way to add any user-defined type to log context.
// Use Object function to pass Loggable as a Field to Logger's methods.
type Loggable interface {
	// ToLog returns a list of fields to log.
	ToLog() []Field
}

// Field is a log context general field.
// Only Logger implementations should interact with Key and Value directly,
// Logger clients should use available Field constructors to fill log context.
type Field struct {
	Key   string
	Value interface{}
}

func Int(key string, value int) Field                { return Field{key, value} }
func Int8(key string, value int8) Field              { return Field{key, value} }
func Int16(key string, value int16) Field            { return Field{key, value} }
func Int32(key string, value int32) Field            { return Field{key, value} }
func Int64(key string, value int64) Field            { return Field{key, value} }
func Uint(key string, value uint) Field              { return Field{key, value} }
func Uint8(key string, value uint8) Field            { return Field{key, value} }
func Uint16(key string, value uint16) Field          { return Field{key, value} }
func Uint32(key string, value uint32) Field          { return Field{key, value} }
func Uint64(key string, value uint64) Field          { return Field{key, value} }
func Float32(key string, value float32) Field        { return Field{key, value} }
func Float64(key string, value float64) Field        { return Field{key, value} }
func Bool(key string, value bool) Field              { return Field{key, value} }
func String(key, value string) Field                 { return Field{key, value} }
func Time(key string, value time.Time) Field         { return Field{key, value} }
func Duration(key string, value time.Duration) Field { return Field{key, value} }
func Error(err error) Field                          { return Field{"error", err} }
func Object(l Loggable) Field                        { return Field{"", l} }

// Nop is a no-op Logger implementation useful in tests.
var Nop Logger = &nop{}

type nop struct{}

func (n *nop) Debug(msg string, fields ...Field) {}
func (n *nop) Info(msg string, fields ...Field)  {}
func (n *nop) Error(msg string, fields ...Field) {}

// callerSkipper is an optional extension for Logger.
// It allows implementations to increase the number of callers skipped by caller annotation.
type callerSkipper interface {
	AddCallerSkip(skip int)
}

// WithFields creates a child Logger that adds the provided fields on each logging operation.
func WithFields(logger Logger, fields ...Field) Logger {
	if skipper, ok := logger.(callerSkipper); ok {
		skipper.AddCallerSkip(1)
	}

	return &withFields{
		logger: logger,
		fields: fields,
	}
}

type withFields struct {
	logger Logger
	fields []Field
}

func (wf *withFields) Debug(msg string, fields ...Field) {
	wf.logger.Debug(msg, append(wf.fields, fields...)...)
}

func (wf *withFields) Info(msg string, fields ...Field) {
	wf.logger.Info(msg, append(wf.fields, fields...)...)
}

func (wf *withFields) Error(msg string, fields ...Field) {
	wf.logger.Error(msg, append(wf.fields, fields...)...)
}

func (wf *withFields) AddCallerSkip(skip int) {
	if skipper, ok := wf.logger.(callerSkipper); ok {
		skipper.AddCallerSkip(skip)
	}
}

// Level indicates a logging priority.
type Level int

const (
	DebugLevel Level = iota - 1
	InfoLevel
	ErrorLevel
)

// Hook is a callback function to be executed before a logging operation.
type Hook func(lvl Level, msg string, fields []Field) error

// WithHooks creates a child Logger that executes the provided hooks on each logging operation.
// If a hook returns an error, it will be logged at ERROR level using the provided logger.
func WithHooks(logger Logger, hooks ...Hook) Logger {
	if skipper, ok := logger.(callerSkipper); ok {
		skipper.AddCallerSkip(1)
	}

	return &withHooks{
		logger: logger,
		hooks:  hooks,
	}
}

type withHooks struct {
	logger Logger
	hooks  []Hook
}

func (wh *withHooks) Debug(msg string, fields ...Field) {
	wh.execHooks(DebugLevel, msg, fields)
	wh.logger.Debug(msg, fields...)
}

func (wh *withHooks) Info(msg string, fields ...Field) {
	wh.execHooks(InfoLevel, msg, fields)
	wh.logger.Info(msg, fields...)
}

func (wh *withHooks) Error(msg string, fields ...Field) {
	wh.execHooks(ErrorLevel, msg, fields)
	wh.logger.Error(msg, fields...)
}

func (wh *withHooks) AddCallerSkip(skip int) {
	if skipper, ok := wh.logger.(callerSkipper); ok {
		skipper.AddCallerSkip(skip)
	}
}

func (wh *withHooks) execHooks(lvl Level, msg string, fields []Field) {
	for _, hook := range wh.hooks {
		if err := hook(lvl, msg, fields); err != nil {
			wh.logger.Error("could not execute hook", Error(err))
		}
	}
}
