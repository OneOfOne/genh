package genh

import (
	"encoding/json"
	"strconv"
	"testing"
)

func TestLMap(t *testing.T) {
	lm := LMapOf(map[string]S{
		"a": {1},
		"b": {2},
		"c": {3},
		"d": {4},
	})

	j, err := json.Marshal(lm)
	DieIf(t, err)

	var lm0 LMap[string, S]
	DieIf(t, json.Unmarshal(j, &lm0))

	if !MapEqual(lm.Raw(), lm0.Raw()) {
		t.Fatal("lm != lm0", lm.Raw(), lm0.Raw())
	}
	j, err = MarshalMsgpack(lm)
	DieIf(t, err)

	var lm1 LMap[string, S]
	DieIf(t, UnmarshalMsgpack(j, &lm1))

	if !MapEqual(lm.Raw(), lm1.Raw()) {
		t.Fatal("lm != lm1", lm.Raw(), lm1.Raw())
	}
}

func TestSLMap(t *testing.T) {
	sm := NewSLMap[S](0)
	for i := 0; i < 10000; i++ {
		sm.Set(strconv.Itoa(i), S{i})
	}
	for _, m := range sm.ms {
		if m.Len() < 250 {
			t.Fatal("m.Len() < 250", m.Len())
		}
	}
	if sm.Len() != 10000 {
		t.Fatal("sm.Len() != 10000", sm.Len())
	}
	count := 0
	sm.ForEach(func(k string, v S) bool {
		count++
		return true
	})
	if count != 10000 {
		t.Fatal("count != 10000", count)
	}
}
