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
	wg.wg.Add(1)
	atomic.AddInt32(&wg.currentConcurrency, 1)
	go wg.worker(f)
}

func (wg *WaitGroup) Wait() {
	wg.wg.Wait()
}

// CurrentConcurrency is useful to inspect runtime concurrency.
func (wg *WaitGroup) CurrentConcurrency() int {
	v := atomic.LoadInt32(&wg.currentConcurrency)
	return int(v)
}

func (wg *WaitGroup) worker(f func()) {
	defer func() {
		atomic.AddInt32(&wg.currentConcurrency, -1)
		wg.wg.Done()
	}()
	f()
}
