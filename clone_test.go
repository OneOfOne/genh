package genh

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

type cloneStruct struct {
	Y map[string]any

	Ptr       *int
	PtrPtr    **int
	PtrPtrPtr ***int
	NilPtr    *int

	S  string
	X  []int
	A  [5]uint64
	x  int
	C  cloner
	C2 *cloner
	C3 cloner0
}

func (c cloneStruct) XV() int {
	return c.x
}

type (
	ifaceBug interface {
		XV() int
	}
	bugSlice []ifaceBug
	cloner   struct {
		A int
	}
)

func (c *cloner) Clone() *cloner {
	return c
}

type cloner0 struct {
	A      int
	cloned bool
}

func (c cloner0) Clone() cloner0 {
	c.cloned = true
	return c
}

func TestBug01(t *testing.T) {
	s := bugSlice{cloneStruct{x: 42}, &cloneStruct{x: 42}}
	c := Clone(s, true)
	t.Log(c[0].XV(), c[1].XV())
}

func TestClone(t *testing.T) {
	n := 42
	pn := &n
	ppn := &pn
	src := &cloneStruct{
		S: "string",
		X: []int{1, 2, 3, 6, 8, 9},
		Y: map[string]any{
			"x": 1, "y": 2.2,
			"z": []int{1, 2, 3, 6, 8, 9},
		},
		Ptr:       pn,
		PtrPtr:    ppn,
		PtrPtrPtr: &ppn,
		A:         [5]uint64{1 << 2, 1 << 4, 1 << 6, 1 << 8, 1 << 10},

		x: n,

		C:  cloner{A: 420},
		C2: &cloner{A: 420},
		C3: cloner0{420, false},
	}

	dst := Clone(src, true)

	if dst == src {
		t.Fatal("cp == s")
	}

	if dst.Ptr == src.Ptr {
		t.Fatal("cp.Ptr == s.Ptr")
	}

	if dst.PtrPtr == src.PtrPtr {
		t.Fatal("cp.PtrPtr == s.PtrPtr")
	}

	if dst.PtrPtrPtr == src.PtrPtrPtr {
		t.Fatal("cp.PtrPtrPtr == s.PtrPtrPtr")
	}

	if src.x != dst.x {
		t.Fatal("src.x != dst.x", src.x, dst.x)
	}

	if !dst.C3.cloned {
		t.Fatal("!dst.C3.cloned")
	}

	dst.C3.cloned = false // so the next check passes

	if !reflect.DeepEqual(src, dst) {
		j1, _ := json.Marshal(src)
		j2, _ := json.Marshal(dst)
		t.Fatalf("!reflect.DeepEqual(src, dst)\nsrc: %s\n----\ndst: %s", j1, j2)
	}

	sj, _ := json.Marshal(src)
	dj, _ := json.Marshal(dst)
	if !bytes.Equal(sj, dj) {
		t.Fatalf("!bytes.Equal(src, dst):\nsrc: %s\ndst: %s", sj, dj)
	}

	dst = Clone(src, false)
	if dst.x == src.x {
		t.Fatal("src.x == dst.x", src.x, dst.x)
	}
	t.Logf("%s", sj)

	if dst.Y["z"].([]int)[0] = 42; src.Y["z"].([]int)[0] != 1 {
		t.Fatal("src.y == dst.y", src.Y, dst.Y)
	}
}

var cloneSink *cloneStruct

func BenchmarkClone(b *testing.B) {
	n := 42
	pn := &n
	ppn := &pn
	s := &cloneStruct{
		S: "string",
		X: []int{1, 2, 3, 6, 8, 9},
		Y: map[string]any{
			"x": 1, "y": 2.2,
			"z": []int{1, 2, 3, 6, 8, 9},
		},
		Ptr:    pn,
		PtrPtr: ppn,
		A:      [5]uint64{1 << 2, 1 << 4, 1 << 6, 1 << 8, 1 << 10},

		x: n,

		C:  cloner{A: 420},
		C2: &cloner{A: 420},
	}
	j, _ := json.Marshal(&s)

	b.Run("Fn", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cloneSink = Clone(s, true)
		}
	})
	b.Run("JSON", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var ss cloneStruct
			if err := json.Unmarshal(j, &ss); err != nil {
				b.Fatal(err)
			}
			cloneSink = &ss
		}
	})
}
