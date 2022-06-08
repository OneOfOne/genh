package genh

import (
	"encoding/json"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

type S struct {
	X int
}

func TestLists(t *testing.T) {
	var l List[S]
	var exp []S
	for i := 0; i < 10; i++ {
		l.Push(S{i})
		exp = append(exp, S{i})
	}

	b, err := json.Marshal(l)
	if err != nil {
		t.Fatal(err)
	}
	var ll List[S]
	err = json.Unmarshal(b, &ll)
	if err != nil {
		t.Fatal(err)
	}
	if !Equal(exp, ll.Slice()) {
		t.Fatal("exp != ll", exp, ll)
	}

	b, err = msgpack.Marshal(l)
	if err != nil {
		t.Fatal(err)
	}
	var ll2 List[S]
	err = msgpack.Unmarshal(b, &ll2)
	if err != nil {
		t.Error(err)
	}
	if !Equal(exp, ll.Slice()) {
		t.Fatal("exp != ll", exp, ll)
	}
}
