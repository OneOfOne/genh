package genh

import (
	"encoding/json"
	"log"
	"math/rand"
	"sort"
	"sync"
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

func TestListClip(t *testing.T) {
	var l, cl, cl2 List[int]
	var nums []int
	for i := 0; i < 25; i++ {
		nums = append(nums, rand.Int()%100)
	}

	for i, n := range nums {
		if i == 10 {
			cl = l.Clip()
		}
		if i == 20 {
			cl2 = l.Clip()
		}
		l.Push(n)
	}

	if l.Len() != 25 || cap(l.Slice()) != l.Len() {
		t.Fatalf("unexpected, should have been 10, got %d %v", cap(l.Slice()), l.Slice())
	}
	cl.Push(99)
	if cl.Len() != 11 || cap(cl.Slice()) != cl.Len() {
		t.Fatalf("unexpected, should have been 11, got %d %v %d", cap(cl.Slice()), cl.Slice(), cl.count())
	}

	if cl2.Len() != 20 || cap(cl2.Slice()) != cl2.Len() {
		t.Fatalf("unexpected, should have been 20, got %d %v %d", cap(cl2.Slice()), cl2.Slice(), cl2.count())
	}

	i := 0
	it := cl2.Iter()
	for v, ok := it.Next(); ok; v, ok = it.Next() {
		if v != nums[i] {
			t.Fatalf("unexpected, should have been %d, got %d", nums[i], v)
		}
		i++
	}
	for i, v := range l.Slice() {
		if v != nums[i] {
			t.Fatalf("unexpected, should have been %d, got %d", nums[i], v)
		}
	}

	// run with race
	var wg sync.WaitGroup
	var mux sync.Mutex
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mux.Lock()
			l.Push(i)
			cl := l.Clip()
			mux.Unlock()
			_ = cl.Slice()
		}(i)
	}
	wg.Wait()
	if l.Len() != 10025 || cap(l.Slice()) != l.Len() {
		t.Fatalf("unexpected, should have been 20, got %d %d", cap(l.Slice()), l.count())
	}
}
