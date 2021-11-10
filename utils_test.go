package log_test

import (
	"reflect"
	"testing"

	"github.com/junk1tm/log"
)

func TestFlattenFields(t *testing.T) {
	tests := []struct {
		name   string
		fields []log.Field
		want   []log.Field
	}{
		{
			name:   "builtin types only",
			fields: []log.Field{log.Int("key_1", 1)},
			want:   []log.Field{log.Int("key_1", 1)},
		},
		{
			name:   "builtin types and Loggable",
			fields: []log.Field{log.Int("key_1", 1), log.Object(A{a: 2})},
			want:   []log.Field{log.Int("key_1", 1), log.Int("key_2", 2)},
		},
		{
			name:   "builtin types and Loggable (nested)",
			fields: []log.Field{log.Int("key_1", 1), log.Object(B{a: A{a: 2}, b: 3})},
			want:   []log.Field{log.Int("key_1", 1), log.Int("key_2", 2), log.Int("key_3", 3)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := log.FlattenFields(tt.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v; want %v", got, tt.want)
			}
		})
	}
}

type A struct {
	a int
}

func (a A) ToLog() []log.Field {
	return []log.Field{
		log.Int("key_2", a.a),
	}
}

type B struct {
	a A
	b int
}

func (b B) ToLog() []log.Field {
	return []log.Field{
		log.Object(b.a),
		log.Int("key_3", b.b),
	}
}
