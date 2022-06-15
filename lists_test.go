package genh

import (
	"encoding/json"
	"testing"
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

	b, err = MarshalMsgpack(l)
	if err != nil {
		t.Fatal(err)
	}
	var ll2 List[S]
	err = UnmarshalMsgpack(b, &ll2)
	if err != nil {
		t.Error(err)
	}
	if !Equal(exp, ll.Slice()) {
		t.Fatal("exp != ll", exp, ll)
	}

	it := l.Iter()
	for v, i := it.Value(), 0; it.Next(); v, i = it.Value(), i+1 {
		t.Log(v, exp[i], i)
		if v != exp[i] {
			t.Fatal("v != exp[i]", v, exp[i])
		}
	}
}
