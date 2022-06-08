package genh

import (
	"encoding/json"
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
