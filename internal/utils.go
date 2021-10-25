package internal

import (
	"github.com/junk1tm/log"
)

// FlattenFields flattens the provided fields slice
// in case it contains Loggable implementations of any nesting,
// thus it's guaranteed that the output is always a slice of builtin types.
// Logger implementations should iterate over FlattenFields(fields)
// to avoid dealing with Loggable directly.
func FlattenFields(fields []log.Field) []log.Field {
	var result []log.Field
	for _, field := range fields {
		if l, ok := field.Value.(log.Loggable); ok {
			result = append(result, FlattenFields(l.Log())...)
		} else {
			result = append(result, field)
		}
	}

	return result
}
