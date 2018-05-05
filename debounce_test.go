package debounce

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func rebuild() {
	fmt.Println("rebuild")
}

func times(n int, dur time.Duration, fun func() error) {
	for i := 0; i < n; i++ {
		fun()
		if dur > 0 {
			time.Sleep(dur)
		}
	}
}

func ExampleNew() {
	db := New(500*time.Millisecond, func() {
		rebuild()
	})

	for i := 0; i < 3; i++ {
		db.Trigger()
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(500 * time.Millisecond)
	// Output: rebuild
	db.Stop()
}

func TestDebounce(t *testing.T) {
	var count int32
	db := New(500*time.Millisecond, func() {
		atomic.AddInt32(&count, 1)
	})
	defer db.Stop()
	times(3, 50*time.Millisecond, db.Trigger)
	time.Sleep(500 * time.Millisecond)

	if atomic.LoadInt32(&count) != 1 {
		t.Errorf("debounce involve error: expect 1, but got %d", atomic.LoadInt32(&count))
	}
}

func TestDebounceStop(t *testing.T) {
	var count int32
	db := New(200*time.Microsecond, func() {
		atomic.AddInt32(&count, 1)
	})
	db.Trigger()
	time.Sleep(500 * time.Microsecond)
	if atomic.LoadInt32(&count) != 1 {
		t.Errorf("debounce involve error: expect 1, but got %d", atomic.LoadInt32(&count))
	}

	db.Stop()

	err := db.Trigger()
	if err != ErrStoped && !db.Stoped() {
		t.Errorf("debounce involve error: trigger a stoped Debouncer should return ErrStoped")
	}

	if atomic.LoadInt32(&count) != 1 {
		t.Errorf("debounce involve error: expect 1, but got %d", atomic.LoadInt32(&count))
	}

	err = db.Stop()
	if err != ErrStoped {
		t.Errorf("debounce involve error: Stop a stoped Debouncer should return ErrStoped")
	}
}
