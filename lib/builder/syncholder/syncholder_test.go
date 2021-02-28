package syncholder

import (
	"testing"
	"time"

	"github.com/avoronkov/go-target-once/lib/targets"
)

func TestSyncHolder(t *testing.T) {
	h := NewResultSyncHolder()

	for i := 0; i < 10; i++ {
		go func() {
			r := h.Get()
			if r.Content != "OK" {
				t.Errorf("Incorrect Get() result: want %v, got %v", "OK", r.Content)
			}
		}()
	}

	time.Sleep(50 * time.Millisecond)
	result := targets.OK("OK")
	h.Put(&result)

	r := h.Get()
	if r.Content != "OK" {
		t.Errorf("Incorrect Get() result: want %v, got %v", "OK", r.Content)
	}

	// Put the value the second time
	res2 := targets.OK("OK2")
	h.Put(&res2)

	for i := 0; i < 10; i++ {
		go func() {
			r2 := h.Get()
			if r2.Content != "OK2" {
				t.Errorf("Incorrect Get() result: want %v, got %v", "OK2", r.Content)
			}
		}()
	}
}
