// Package debounce implement Golang version debounce.
// The passed function which will postpone its execution until after `time.Duration`
// have elapsed since the last time it was invoked.
package debounce

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// ErrStoped raised by Trigger or Stop a stoped Debouncer
var ErrStoped = errors.New("debounce stoped")

// Debouncer define debounce interfac
type Debouncer interface {
	// Stop/Cancel the Debounce. Stop a stoped Debouncer will return ErrStoped
	Stop() error
	// Trigger a invocation. Trigger a stoped Debouncer will return ErrStoped
	Trigger() error
	// Stoped checks if Debouncer is stoped
	Stoped() bool
}

type debounced struct {
	// lazy create goroutine
	runner sync.Once
	period time.Duration
	stoped int32
	// trigger signal
	in chan struct{}
	// involve signal
	out chan struct{}
	// done
	done chan struct{}
}

func (db *debounced) Trigger() error {
	if atomic.LoadInt32(&db.stoped) == 1 {
		return ErrStoped
	}

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
	return nil
}

func (db *debounced) Stop() error {
	if !atomic.CompareAndSwapInt32(&db.stoped, 0, 1) {
		return ErrStoped
	}
	db.done <- struct{}{}
	close(db.in)
	close(db.out)
	close(db.done)
	return nil
}

func (db *debounced) Stoped() bool {
	return atomic.LoadInt32(&db.stoped) == 1
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
