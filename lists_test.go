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
	for i := 0; i < 10; i++ {
		l.Push(S{i})
		t.Log(l.Slice())
	}

	b, err := json.Marshal(l)
	t.Logf("%s: %v", b, err)
	var ll List[S]
	err = json.Unmarshal(b, &ll)
	t.Logf("%v: %v", err, ll.Slice())
}
