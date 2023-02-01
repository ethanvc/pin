package conc

import (
	"sync"
	"sync/atomic"
)

type Pool struct {
	maxConcurrencyCount     int32
	currentConcurrencyCount int32
	waitG                   sync.WaitGroup
	initGuard               sync.Once
	task                    chan func()
	manageTask              chan struct{}
}

func (p *Pool) SetMaxConcurrencyCount(count int) {
	p.init()
	atomic.StoreInt32(&p.maxConcurrencyCount, int32(count))
	select {
	case p.manageTask <- struct{}{}:
	default:
	}
}

func (p *Pool) GetMaxConcurrencyCount() int {
	v := atomic.LoadInt32(&p.maxConcurrencyCount)
	return int(v)
}

func (p *Pool) GetCurrentConcurrencyCount() int {
	v := atomic.LoadInt32(&p.currentConcurrencyCount)
	return int(v)
}

func (p *Pool) Go(f func()) {
	p.init()
	p.task <- f
}

func (p *Pool) Cleanup() {
	p.init()
	close(p.manageTask)
	p.waitG.Wait()
}

func (p *Pool) init() {
	p.initGuard.Do(func() {
		p.task = make(chan func())
		p.manageTask = make(chan struct{})
		// pump mange routine working
		p.manageTask <- struct{}{}
		go p.manageRoutine()
	})
}

func (p *Pool) needClose() bool {
	select {
	case _, ok := <-p.manageTask:
		return !ok
	default:
	}
	return false
}

func (p *Pool) manageRoutine() {
	for {
		if p.needClose() {
			break
		}
		p.parkHere()

		select {
		case <-p.manageTask:
			continue
		case f := <-p.task:
			p.dispatchWork(f)
		}
	}
}

func (p *Pool) parkHere() {
	for {
		if atomic.LoadInt32(&p.maxConcurrencyCount) == 0 {
			<-p.manageTask
			continue
		}
		break
	}
}

func (p *Pool) dispatchWork(f func()) {
	p.addConcurrency()
	// plus one because manage routine can execute work too.
	if atomic.LoadInt32(&p.currentConcurrencyCount)+1 < atomic.LoadInt32(&p.maxConcurrencyCount) {
		go p.workRoutine(f)
	} else {
		p.workRoutine(f)
	}
}

func (p *Pool) workRoutine(f func()) {
	defer p.subConcurrency()
	f()
	for {
		if p.workerNeedExit() {
			break
		}
		select {
		case f = <-p.task:
			f()
		default:
			break
		}
	}
}

func (p *Pool) workerNeedExit() bool {
	return atomic.LoadInt32(&p.currentConcurrencyCount) > atomic.LoadInt32(&p.maxConcurrencyCount)
}

func (p *Pool) addConcurrency() {
	p.waitG.Add(1)
	atomic.AddInt32(&p.currentConcurrencyCount, 1)
}

func (p *Pool) subConcurrency() {
	atomic.AddInt32(&p.currentConcurrencyCount, -1)
	p.waitG.Done()
}
