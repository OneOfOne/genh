package genh

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestTimedMap(t *testing.T) {
	n := runtime.NumGoroutine()
	var tm TimedMap[string, string]
	i := 0
	tm.SetUpdateExpireFn("key", func() string {
		i++
		return fmt.Sprintf("val:%d", i)
	}, time.Millisecond*110, time.Millisecond*400)
	time.Sleep(time.Millisecond * 100)
	if v := tm.Get("key"); v != "val:1" {
		t.Fatal("expected val:1, got", v)
	}
	time.Sleep(time.Millisecond * 350)
	if v := tm.Get("key"); v != "val:5" {
		t.Fatal("expected val:5, got", v)
	}
	time.Sleep(time.Millisecond * 500)
	if v := tm.Get("key"); v != "" {
		t.Fatal("expected empty, got", v)
	}

	tm.Set("key2", "val1", time.Millisecond*100)
	time.Sleep(time.Millisecond * 350)
	if v := tm.Get("key2"); v != "" {
		t.Fatal("expected empty, got", v)
	}
	tm.Set("key2", "val2", time.Minute)
	time.Sleep(time.Millisecond * 350)
	tm.Set("key2", "val3", time.Minute)

	if v, ok := tm.GetOk("key2"); v != "val3" || !ok {
		t.Fatalf("unexpected val3: %v", v)
	}
	tm.Delete("key2")

	if nn := runtime.NumGoroutine(); nn != n {
		t.Fatalf("goroutine leak: %d > %d", nn, n)
	}
}
