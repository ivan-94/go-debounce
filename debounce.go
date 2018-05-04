// Package debounce implement Golang version debounce.
// The passed function which will postpone its execution until after `time.Duration`
// have elapsed since the last time it was invoked.
package debounce

import (
	"sync"
	"time"
)

// Debouncer define debounce interfac
type Debouncer interface {
	// Stop/Cancel the Debounce
	Stop()
	// Trigger a invocation
	Trigger()
}

type debounced struct {
	// lazy create goroutine
	runner sync.Once
	period time.Duration
	// trigger signal
	in chan struct{}
	// involve signal
	out chan struct{}
	// done
	done chan struct{}
}

func (db *debounced) Trigger() {
	db.runner.Do(func() {
		go func() {
			for {
				select {
				case <-db.in:
					// do nothing, just reset timer
				case <-time.After(db.period):
					// timeouted, involve callback
					db.out <- struct{}{}
					// wait for next trigger
					<-db.in
				case <-db.done:
					return
				}
			}
		}()
	})
	db.in <- struct{}{}
}

func (db *debounced) Stop() {
	db.done <- struct{}{}
	close(db.in)
	close(db.out)
	close(db.done)
}

// New create a debouncer
func New(wait time.Duration, callback func()) Debouncer {
	db := &debounced{
		period: wait,
		in:     make(chan struct{}),
		out:    make(chan struct{}),
		done:   make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-db.out:
				callback()
			case <-db.done:
				return
			}
		}
	}()

	return db
}
