// Package zerologimpl contains zerolog implementation of Logger interface.
package zerologimpl

import (
	"fmt"
	"time"

	"github.com/junk1tm/log"
	"github.com/rs/zerolog"
)

// NewLogger creates a new log.Logger from the provided zerolog.Logger.
func NewLogger(logger zerolog.Logger) log.Logger {
	return &wrapper{
		callerSkip: 1,
		logger:     logger,
	}
}

type wrapper struct {
	callerSkip int
	logger     zerolog.Logger
}

func (w *wrapper) Debug(msg string, fields ...log.Field) { w.log(w.logger.Debug(), msg, fields) }
func (w *wrapper) Info(msg string, fields ...log.Field)  { w.log(w.logger.Info(), msg, fields) }
func (w *wrapper) Error(msg string, fields ...log.Field) { w.log(w.logger.Error(), msg, fields) }

func (w *wrapper) AddCallerSkip(skip int) { w.callerSkip += skip }

func (w *wrapper) log(event *zerolog.Event, msg string, fields []log.Field) {
	for _, field := range log.FlattenFields(fields) {
		switch value := field.Value.(type) {
		case int:
			event.Int(field.Key, value)
		case int8:
			event.Int8(field.Key, value)
		case int16:
			event.Int16(field.Key, value)
		case int32:
			event.Int32(field.Key, value)
		case int64:
			event.Int64(field.Key, value)
		case uint:
			event.Uint(field.Key, value)
		case uint8:
			event.Uint8(field.Key, value)
		case uint16:
			event.Uint16(field.Key, value)
		case uint32:
			event.Uint32(field.Key, value)
		case uint64:
			event.Uint64(field.Key, value)
		case float32:
			event.Float32(field.Key, value)
		case float64:
			event.Float64(field.Key, value)
		case bool:
			event.Bool(field.Key, value)
		case string:
			event.Str(field.Key, value)
		case time.Time:
			event.Time(field.Key, value)
		case time.Duration:
			event.Dur(field.Key, value)
		case error:
			event.AnErr(field.Key, value)
		default:
			panic(fmt.Sprintf("unexpected field type %T", value))
		}
	}

	event.CallerSkipFrame(w.callerSkip + 1).Msg(msg)
}
