package genh

import (
	"reflect"
	"testing"
)

func TestFirstNonZero(t *testing.T) {
	type args struct {
		args []any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{"ints", args{[]any{nil, 1, 2, 3}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FirstNonZero(tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FirstNonZero() = %v, want %v", got, tt.want)
			}
		})
	}
}
