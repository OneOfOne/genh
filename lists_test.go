package genh

import (
	"encoding/json"
	"log"
	"math/rand"
	"sort"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

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

	lls := l.ListAt(5, 10)
	if !Equal(lls.Slice(), exp[5:]) || lls.Len() != 5 {
		t.Fatal("exp != ll", exp[5:], lls.Slice(), lls.Len())
	}
	lls = l.ListAt(5, -1)
	if !Equal(lls.Slice(), exp[5:]) || lls.Len() != 5 {
		t.Fatal("exp != ll", exp[5:], lls.Slice(), lls.Len())
	}
}

func TestListIter(t *testing.T) {
	var l List[S]
	var exp []S
	for i := 5; i < 10; i++ {
		l.Push(S{i})
		exp = append(exp, S{i})
	}
	it := l.Iter()
	i := 0
	for v, ok := it.Next(); ok; v, ok = it.Next() {
		if v != exp[i] {
			t.Fatal("v != exp[i]", v, exp[i])
		}
		i++
	}

	l.Prepend(S{44})
	it = l.Iter()
	for v, ok := it.Next(); ok; v, ok = it.Next() {
		if v.X == 7 {
			it.Set(S{77})
		}
		if v.X == 5 || v.X == 8 {
			it.Delete()
		}

	}

	exp = []S{{44}, {6}, {77}, {9}}
	if !Equal(l.Slice(), exp) || l.Len() != 4 {
		t.Fatal("exp != ll", exp, l.Slice(), l.Len())
	}
}

func TestListSort(t *testing.T) {
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
