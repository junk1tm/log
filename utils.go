package log

// FlattenFields flattens the provided fields slice
// in case it contains Loggable implementations of any nesting,
// thus it's guaranteed that the output is always a slice of builtin types.
// This function is intended to be used by Logger implementations
// that should iterate over FlattenFields(fields) to avoid dealing with Loggable directly.
func FlattenFields(fields []Field) []Field {
	var result []Field
	for _, field := range fields {
		if l, ok := field.Value.(Loggable); ok {
			result = append(result, FlattenFields(l.ToLog())...)
		} else {
			result = append(result, field)
		}
	}

	return result
}
