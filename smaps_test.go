package genh

import (
	"runtime"
	"strconv"
	"sync"
	"testing"
)

func BenchmarkLMaps(b *testing.B) {
	N := runtime.NumCPU() * 100
	var keys [100]string
	for i := range keys {
		keys[i] = strconv.Itoa(i)
	}
	b.Run("LMap", func(b *testing.B) {
		var m LMap[string, int]
		for b.Loop() {
			var wg sync.WaitGroup
			for j := range N {
				wg.Add(1)
				j := j % len(keys)
				go func() {
					if m.MustGet(keys[j], func() int { return j }) != j {
						panic("bad j")
					}
					wg.Done()
				}()
			}
			wg.Wait()
		}
	})
	b.Run("SLMap", func(b *testing.B) {
		var m SLMap[int]
		for b.Loop() {
			var wg sync.WaitGroup
			for j := range N {
				wg.Add(1)
				j := j % len(keys)
				go func() {
					if m.MustGet(keys[j], func() int { return j }) != j {
						panic("bad j")
					}
					wg.Done()
				}()
			}
			wg.Wait()
		}
	})
}

func BenchmarkLMultiMaps(b *testing.B) {
	N := runtime.NumCPU() * 100
	var keys [100]string
	for i := range keys {
		keys[i] = strconv.Itoa(i)
	}
	b.Run("LMultiMap", func(b *testing.B) {
		var m LMultiMap[string, string, int]
		for b.Loop() {
			var wg sync.WaitGroup
			for j := range N {
				wg.Add(1)
				j := j % len(keys)
				go func() {
					if m.MustGet(keys[j], keys[j], func() int { return j }) != j {
						panic("bad j")
					}
					wg.Done()
				}()
			}
			wg.Wait()
		}
	})
	b.Run("SLMultiMap", func(b *testing.B) {
		var m SLMultiMap[int]
		for b.Loop() {
			var wg sync.WaitGroup
			for j := range N {
				wg.Add(1)
				j := j % len(keys)
				go func() {
					if m.MustGet(keys[j], keys[j], func() int { return j }) != j {
						panic("bad j")
					}
					wg.Done()
				}()
			}
			wg.Wait()
		}
	})
}
