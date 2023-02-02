package conc

import (
	"sync"
	"sync/atomic"
)

const needCleanupConcurrencyCount = -1

type Pool struct {
	maxConcurrencyCount     int32
	currentConcurrencyCount int32
	waitG                   sync.WaitGroup
	initGuard               sync.Once
	task                    chan func()
	manageTask              chan struct{}
	quitChan                chan struct{}
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
	v := atomic.LoadInt32(&p.currentConcurrencyCount)
	return int(v)
}

func (p *Pool) Go(f func()) {
	p.init()
	p.task <- f
}

func (p *Pool) Wait() {
	p.waitG.Wait()
}

func (p *Pool) Cleanup() {
	p.init()
	p.quitChan = make(chan struct{})
	close(p.task)
	p.setMaxConcurrencyCount(needCleanupConcurrencyCount)

	<-p.quitChan
	p.waitG.Wait()
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

		for {
			select {
			case <-p.manageTask:
				break
			case f, ok := <-p.task:
				if !ok {
					break
				}
				p.dispatchWork(f)
			}
		}
	}
	close(p.quitChan)
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
	p.addConcurrency()
	// plus one because manage routine can execute work too.
	if atomic.LoadInt32(&p.currentConcurrencyCount)+1 < atomic.LoadInt32(&p.maxConcurrencyCount) {
		go p.workRoutine(f)
	} else {
		defer p.subConcurrency()
		f()
	}
}

func (p *Pool) workRoutine(f func()) {
	defer p.subConcurrency()
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
