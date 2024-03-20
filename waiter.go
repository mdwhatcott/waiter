package waiter

import (
	"errors"
	"sync/atomic"
)

type Waiter interface {
	Add(delta int)
	Done()
	Wait()
}

type waiter struct {
	counter *atomic.Int64
	waits   *atomic.Int64
}

func New() Waiter {
	return &waiter{
		counter: new(atomic.Int64),
		waits:   new(atomic.Int64),
	}
}
func (this *waiter) Add(delta int) {
	if this.waits.Load() > 0 {
		panic(ErrAddDuringWait)
	}
	this.counter.Add(int64(delta))
}
func (this *waiter) Done() {
	counter := this.counter.Add(-1)
	if counter < 0 {
		panic(ErrNegativeCounter)
	}
}
func (this *waiter) Wait() {
	this.waits.Add(1)
	defer this.waits.Add(-1)
	for {
		value := this.counter.Load()
		if value < 0 {
			panic(ErrNegativeCounter)
		} else if value == 0 {
			break
		}
	}
}

var (
	ErrNegativeCounter = errors.New("negative counter")
	ErrAddDuringWait   = errors.New("cannot add while waiting")
)
