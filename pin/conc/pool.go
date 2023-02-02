package conc

import (
	"sync"
	"sync/atomic"
)

const needCleanupConcurrencyCount = -1

type Pool struct {
	maxConcurrencyCount int32
	wg                  WaitGroup
	initGuard           sync.Once
	task                chan func()
	manageTask          chan struct{}
	quitChan            atomic.Pointer[chan struct{}]
}

func (p *Pool) SetMaxConcurrencyCount(count int) {
	if count < 0 {
		return
	}
	p.setMaxConcurrencyCount(count)
}

func (p *Pool) setMaxConcurrencyCount(count int) {
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
	return p.wg.CurrentConcurrency()
}

func (p *Pool) Go(f func()) {
	p.init()
	p.task <- f
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Cleanup() {
	p.init()
	qc := make(chan struct{})
	p.quitChan.Store(&qc)
	close(p.task)
	p.setMaxConcurrencyCount(needCleanupConcurrencyCount)

	<-qc
	p.wg.Wait()
}

func (p *Pool) init() {
	p.initGuard.Do(func() {
		p.task = make(chan func())
		p.manageTask = make(chan struct{}, 1)
		// pump mange routine working
		p.manageTask <- struct{}{}
		go p.manageRoutine()
	})
}

func (p *Pool) manageRoutine() {
	for {
		if p.parkHere() {
			break
		}

	Exit:
		for {
			select {
			case <-p.manageTask:
				break Exit
			case f, ok := <-p.task:
				if !ok {
					break Exit
				}
				p.dispatchWork(f)
			}
		}
	}
	close(*p.quitChan.Load())
}

func (p *Pool) parkHere() bool {
	for {
		cnt := atomic.LoadInt32(&p.maxConcurrencyCount)
		if cnt == needCleanupConcurrencyCount {
			return true
		}
		if cnt == 0 {
			<-p.manageTask
			continue
		}
		break
	}
	return false
}

func (p *Pool) dispatchWork(f func()) {
	p.wg.Add(1)
	// plus one because manage routine can execute work too.
	if p.wg.CurrentConcurrency()+1 < p.GetMaxConcurrencyCount() {
		go p.workRoutine(f)
	} else {
		defer p.wg.Done()
		f()
	}
}

func (p *Pool) workRoutine(f func()) {
	defer p.wg.Done()
	f()
	for {
		if p.WorkerNeedExit() {
			break
		}
		select {
		case f, ok := <-p.task:
			if !ok {
				break
			}
			f()
		default:
			break
		}
	}
}

func (p *Pool) WorkerNeedExit() bool {
	return p.GetCurrentConcurrencyCount() >= p.GetMaxConcurrencyCount()
}
