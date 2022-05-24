package genh

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

type cloneStruct struct {
	Y      map[string]any
	Ptr    *int
	NilPtr *int
	S      string
	X      []int
	A      [5]uint64
}

func TestTypedClone(t *testing.T) {
	n := 42
	src := &cloneStruct{
		S: "string",
		X: []int{1, 2, 3, 6, 8, 9},
		Y: map[string]any{
			"x": 1, "y": 2.2,
			"z": []int{1, 2, 3, 6, 8, 9},
		},
		Ptr: &n,
		A:   [5]uint64{1 << 2, 1 << 4, 1 << 6, 1 << 8, 1 << 10},
	}

	dst := TypedClone(src)

	if dst == src {
		t.Fatal("cp == s")
	}

	if dst.Ptr == src.Ptr {
		t.Fatal("cp.Ptr == s.Ptr")
	}

	if !reflect.DeepEqual(src, dst) {
		t.Fatal("!reflect.DeepEqual(src, dst)")
	}

	sj, _ := json.Marshal(src)
	dj, _ := json.Marshal(dst)
	if !bytes.Equal(sj, dj) {
		t.Fatalf("!bytes.Equal(src, dst):\nsrc: %s\ndst: %s", sj, dj)
	}
	t.Logf("%s", sj)
}

func BenchmarkTypedClone(b *testing.B) {
	s := cloneStruct{
		S: "string",
		X: []int{1, 2, 3, 6, 8, 9},
		Y: map[string]any{
			"x": 1, "y": 2.2,
			"z": []int{1, 2, 3, 6, 8, 9},
		},
		A: [5]uint64{1 << 2, 1 << 4, 1 << 6, 1 << 8, 1 << 10},
	}
	j, _ := json.Marshal(&s)

	b.Run("Fn", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if !reflect.DeepEqual(&s, TypedClone(&s)) {
				b.Fatal("bad")
			}
		}
	})
	b.Run("JSON", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var ss cloneStruct
			if err := json.Unmarshal(j, &ss); err != nil {
				b.Fatal(err)
			}
			jj, _ := json.Marshal(&ss)
			if !bytes.Equal(jj, j) {
				b.Fatalf("bad\n%s\n%s", j, jj)
			}
		}
	})
}
