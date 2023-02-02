package conc

import (
	"sync"
	"sync/atomic"
)

type WaitGroup struct {
	wg                 sync.WaitGroup
	currentConcurrency int32
}

func (wg *WaitGroup) Go(f func()) {
	wg.Add(1)
	go wg.worker(f)
}

func (wg *WaitGroup) Wait() {
	wg.wg.Wait()
}

func (wg *WaitGroup) Add(delta int) {
	wg.wg.Add(delta)
	atomic.AddInt32(&wg.currentConcurrency, int32(delta))
}

func (wg *WaitGroup) Done() {
	atomic.AddInt32(&wg.currentConcurrency, -1)
	wg.wg.Done()
}

// CurrentConcurrency is useful to inspect runtime concurrency.
func (wg *WaitGroup) CurrentConcurrency() int {
	v := atomic.LoadInt32(&wg.currentConcurrency)
	return int(v)
}

func (wg *WaitGroup) worker(f func()) {
	defer wg.Done()
	f()
}
