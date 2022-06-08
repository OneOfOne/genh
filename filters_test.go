package genh

import (
	"reflect"
	"sort"
	"testing"
)

func TestGroupBy(t *testing.T) {
	type args struct {
		in map[string]int
		fn func(k string, v int) string
	}
	tests := []struct {
		name string
		args args
		want map[string][]int
	}{
		{
			name: "odds and evens",
			args: args{
				in: map[string]int{
					"1":  1,
					"2":  2,
					"3":  3,
					"4":  4,
					"5":  5,
					"6":  6,
					"7":  7,
					"8":  8,
					"9":  9,
					"10": 10,
				},
				fn: func(k string, v int) string {
					if v != 0 && v%2 == 0 {
						return "even"
					}
					return "odd"
				},
			},
			want: map[string][]int{
				"odd":  {1, 3, 5, 7, 9},
				"even": {2, 4, 6, 8, 10},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortMap(GroupBy(tt.args.in, tt.args.fn)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func sortMap(m map[string][]int) map[string][]int {
	for _, v := range m {
		sort.Ints(v)
	}
	return m
}
