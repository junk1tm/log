// Package zap contains zap implementation of Logger interface.
package zap

import (
	"fmt"
	"time"

	"github.com/junk1tm/log"
	"github.com/junk1tm/log/impl"
	"go.uber.org/zap"
)

// NewLogger creates a new log.Logger from the provided zap.Logger.
func NewLogger(logger *zap.Logger) log.Logger {
	return &wrapper{
		logger: logger.WithOptions(zap.AddCallerSkip(1)),
	}
}

type wrapper struct {
	logger *zap.Logger
}

func (w *wrapper) Debug(msg string, fields ...log.Field) { w.logger.Debug(msg, zapFields(fields)...) }
func (w *wrapper) Info(msg string, fields ...log.Field)  { w.logger.Info(msg, zapFields(fields)...) }
func (w *wrapper) Error(msg string, fields ...log.Field) { w.logger.Error(msg, zapFields(fields)...) }

func (w *wrapper) AddCallerSkip(skip int) { w.logger = w.logger.WithOptions(zap.AddCallerSkip(skip)) }

func zapFields(fields []log.Field) []zap.Field {
	var zf []zap.Field

	for _, field := range impl.FlattenFields(fields) {
		switch value := field.Value.(type) {
		case int:
			zf = append(zf, zap.Int(field.Key, value))
		case int8:
			zf = append(zf, zap.Int8(field.Key, value))
		case int16:
			zf = append(zf, zap.Int16(field.Key, value))
		case int32:
			zf = append(zf, zap.Int32(field.Key, value))
		case int64:
			zf = append(zf, zap.Int64(field.Key, value))
		case uint:
			zf = append(zf, zap.Uint(field.Key, value))
		case uint8:
			zf = append(zf, zap.Uint8(field.Key, value))
		case uint16:
			zf = append(zf, zap.Uint16(field.Key, value))
		case uint32:
			zf = append(zf, zap.Uint32(field.Key, value))
		case uint64:
			zf = append(zf, zap.Uint64(field.Key, value))
		case float32:
			zf = append(zf, zap.Float32(field.Key, value))
		case float64:
			zf = append(zf, zap.Float64(field.Key, value))
		case bool:
			zf = append(zf, zap.Bool(field.Key, value))
		case string:
			zf = append(zf, zap.String(field.Key, value))
		case time.Time:
			zf = append(zf, zap.Time(field.Key, value))
		case time.Duration:
			zf = append(zf, zap.Duration(field.Key, value))
		case error:
			zf = append(zf, zap.NamedError(field.Key, value))
		default:
			panic(fmt.Sprintf("unexpected field type %T", value))
		}
	}

	return zf
}
