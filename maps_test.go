// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package genh

import (
	"math"
	"sort"
	"testing"
)

var (
	m1 = map[int]int{1: 1, 2: 4, 3: 3, 4: 8, 5: 5, 6: 12, 7: 7, 8: 16}
	m2 = map[int]string{1: "2", 2: "4", 4: "8", 8: "16"}
)

func TestMapKeys(t *testing.T) {
	want := []int{1, 2, 3, 4, 5, 6, 7, 8}

	got1 := MapKeys(m1)
	sort.Ints(got1)
	if !Equal(got1, want) {
		t.Errorf("Keys(%v) = %v, want %v", m1, got1, want)
	}
}

func TestMapValues(t *testing.T) {
	got1 := MapValues(m1)
	want1 := []int{1, 3, 4, 5, 7, 8, 12, 16}
	sort.Ints(got1)
	if !Equal(got1, want1) {
		t.Errorf("Values(%v) = %v, want %v", m1, got1, want1)
	}

	got2 := MapValues(m2)
	want2 := []string{"16", "2", "4", "8"}
	sort.Strings(got2)
	if !Equal(got2, want2) {
		t.Errorf("Values(%v) = %v, want %v", m2, got2, want2)
	}
}

func TestMapEqual(t *testing.T) {
	if !MapEqual(m1, m1) {
		t.Errorf("MapEqual(%v, %v) = false, want true", m1, m1)
	}
	if MapEqual(m1, (map[int]int)(nil)) {
		t.Errorf("MapEqual(%v, nil) = true, want false", m1)
	}
	if MapEqual((map[int]int)(nil), m1) {
		t.Errorf("MapEqual(nil, %v) = true, want false", m1)
	}
	if !MapEqual[map[int]int, map[int]int](nil, nil) {
		t.Error("MapEqual(nil, nil) = false, want true")
	}
	if ms := map[int]int{1: 2}; MapEqual(m1, ms) {
		t.Errorf("MapEqual(%v, %v) = true, want false", m1, ms)
	}

	// Comparing NaN for equality is expected to fail.
	mf := map[int]float64{1: 0, 2: math.NaN()}
	if MapEqual(mf, mf) {
		t.Errorf("MapEqual(%v, %v) = true, want false", mf, mf)
	}
}

func TestMapEqualFunc(t *testing.T) {
	if !MapEqualFunc(m1, m1, equal[int]) {
		t.Errorf("MapEqualFunc(%v, %v, equal) = false, want true", m1, m1)
	}
	if MapEqualFunc(m1, (map[int]int)(nil), equal[int]) {
		t.Errorf("MapEqualFunc(%v, nil, equal) = true, want false", m1)
	}
	if MapEqualFunc((map[int]int)(nil), m1, equal[int]) {
		t.Errorf("MapEqualFunc(nil, %v, equal) = true, want false", m1)
	}
	if !MapEqualFunc[map[int]int, map[int]int](nil, nil, equal[int]) {
		t.Error("MapEqualFunc(nil, nil, equal) = false, want true")
	}
	if ms := map[int]int{1: 2}; MapEqualFunc(m1, ms, equal[int]) {
		t.Errorf("MapEqualFunc(%v, %v, equal) = true, want false", m1, ms)
	}

	// Comparing NaN for equality is expected to fail.
	mf := map[int]float64{1: 0, 2: math.NaN()}
	if MapEqualFunc(mf, mf, equal[float64]) {
		t.Errorf("MapEqualFunc(%v, %v, equal) = true, want false", mf, mf)
	}
	// But it should succeed using equalNaN.
	if !MapEqualFunc(mf, mf, equalNaN[float64]) {
		t.Errorf("MapEqualFunc(%v, %v, equalNaN) = false, want true", mf, mf)
	}

	// if !MapEqualFunc(m1, m2, equalIntStr) {
	// 	t.Errorf("MapEqualFunc(%v, %v, equalIntStr) = false, want true", m1, m2)
	// }
}

func TestMapClear(t *testing.T) {
	ml := map[int]int{1: 1, 2: 2, 3: 3}
	MapClear(ml)
	if got := len(ml); got != 0 {
		t.Errorf("len(%v) = %d after Clear, want 0", ml, got)
	}
	if !MapEqual(ml, (map[int]int)(nil)) {
		t.Errorf("MapEqual(%v, nil) = false, want true", ml)
	}
}

func TestMapClone(t *testing.T) {
	mc := MapClone(m1)
	if !MapEqual(mc, m1) {
		t.Errorf("Clone(%v) = %v, want %v", m1, mc, m1)
	}
	mc[16] = 32
	if MapEqual(mc, m1) {
		t.Errorf("MapEqual(%v, %v) = true, want false", mc, m1)
	}
}

func TestMapCopy(t *testing.T) {
	mc := MapClone(m1)
	MapCopy(mc, mc)
	if !MapEqual(mc, m1) {
		t.Errorf("Copy(%v, %v) = %v, want %v", m1, m1, mc, m1)
	}
	MapCopy(mc, map[int]int{16: 32})
	want := map[int]int{1: 1, 2: 4, 3: 3, 4: 8, 5: 5, 6: 12, 7: 7, 8: 16, 16: 32}
	if !MapEqual(mc, want) {
		t.Errorf("Copy result = %v, want %v", mc, want)
	}
}

func TestMapDeleteFunc(t *testing.T) {
	mc := MapClone(m1)
	MapDeleteFunc(mc, func(int, int) bool { return false })
	if !MapEqual(mc, m1) {
		t.Errorf("DeleteFunc(%v, true) = %v, want %v", m1, mc, m1)
	}
	MapDeleteFunc(mc, func(k, v int) bool { return k > 3 })
	want := map[int]int{1: 1, 2: 4, 3: 3}
	if !MapEqual(mc, want) {
		t.Errorf("DeleteFunc result = %v, want %v", mc, want)
	}
}

func TestMapFilter(t *testing.T) {
	mc := MapClone(m1)
	mc = MapFilter(mc, func(_ int, v int) bool { return v%2 == 0 }, false)
	for _, v := range mc {
		if v%2 != 0 {
			t.Fatalf("bad v %d", v)
		}
	}
}
