// Package log provides a general-purpose logging API.
package log

import (
	"time"
)

// Logger is a minimal logging interface.
// It uses structured logging with Field parameters to ensure type safety.
type Logger interface {
	Debug(text string, fields ...Field)
	Info(text string, fields ...Field)
	Error(text string, fields ...Field)
}

// Loggable represents an entity that can be logged.
// It provides a type-safe way to add any user-defined type to log context.
// Use Object function to pass Loggable as a Field to Logger's methods.
type Loggable interface {
	// Log returns a list of fields to log.
	Log() []Field
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
func String(key string, value string) Field          { return Field{key, value} }
func Time(key string, value time.Time) Field         { return Field{key, value} }
func Duration(key string, value time.Duration) Field { return Field{key, value} }
func Error(err error) Field                          { return Field{"error", err} }
func Object(l Loggable) Field                        { return Field{"", flatten(l)} }

func flatten(l Loggable) []Field {
	var fields []Field
	for _, field := range l.Log() {
		if fs, ok := field.Value.([]Field); ok {
			fields = append(fields, fs...)
		} else {
			fields = append(fields, field)
		}
	}
	return fields
}

// Nop is a no-op Logger implementation useful in tests.
var Nop Logger = &nop{}

type nop struct{}

func (n *nop) Debug(text string, fields ...Field) {}
func (n *nop) Info(text string, fields ...Field)  {}
func (n *nop) Error(text string, fields ...Field) {}
