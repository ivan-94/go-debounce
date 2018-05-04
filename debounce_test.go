package debounce

import (
	"fmt"
	"time"
)

func rebuild() {
	fmt.Println("rebuild")
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
