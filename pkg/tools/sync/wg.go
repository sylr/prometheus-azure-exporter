package sync

// Orininal code from: https://play.golang.org/p/ChjP2wpvyt

import (
	"sync"
)

// Waiter interface
type Waiter interface {
	Add(int)
	Done()
	Wait()
}

// NewWaiter returns a Waiter, a sync.WaitGroup if cap i <= 0, a BoundedWaitGroup otherwise.
func NewWaiter(cap int) Waiter {
	if cap > 0 {
		wg := NewBoundedWaitGroup(cap)
		return &wg
	} else {
		return &sync.WaitGroup{}
	}
}

// BoundedWaitGroup is a wait group which has a limit boundary meaning it will
// wait for Done() to be called before releasing Add(n) if the limit has been reached
type BoundedWaitGroup struct {
	wg sync.WaitGroup
	ch chan struct{}
}

// NewBoundedWaitGroup returns a new BoundedWaitGroup
func NewBoundedWaitGroup(cap int) BoundedWaitGroup {
	return BoundedWaitGroup{ch: make(chan struct{}, cap)}
}

// Add ...
func (bwg *BoundedWaitGroup) Add(delta int) {
	// Delta < 0
	for i := 0; i > delta; i-- {
		<-bwg.ch
	}

	// Delta > 0
	for i := 0; i < delta; i++ {
		bwg.ch <- struct{}{}
	}
	bwg.wg.Add(delta)
}

// Done ...
func (bwg *BoundedWaitGroup) Done() {
	bwg.Add(-1)
}

// Wait ...
func (bwg *BoundedWaitGroup) Wait() {
	bwg.wg.Wait()
}
