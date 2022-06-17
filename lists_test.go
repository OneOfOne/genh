package genh

import (
	"encoding/json"
	"log"
	"math/rand"
	"sort"
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
	i := 0
	for v, ok := it.Next(); ok; v, ok = it.Next() {
		if v != exp[i] {
			t.Fatal("v != exp[i]", v, exp[i])
		}
		i++
	}

	lls := l.ListAt(5, -1)
	if !Equal(lls.Slice(), exp[5:]) || lls.Len() != 5 {
		t.Fatal("exp != ll", exp[5:], lls, lls.Len())
	}
}

func TestListSort(t *testing.T) {
	log.SetFlags(0)
	var l List[int]
	var nums []int
	for i := 0; i < 25; i++ {
		nums = append(nums, rand.Int()%100)
	}

	for _, n := range nums {
		l.PushSort(n, func(a, b int) bool { return a >= b })
	}

	sort.Sort(sort.Reverse(sort.IntSlice(nums)))
	if !Equal(nums, l.Slice()) {
		t.Log(nums)
		t.Log(l.Slice())
		t.Fatal("neq")
	}

	t.Log(l.Slice())
}
