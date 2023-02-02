package conc

import (
	"sync"
	"sync/atomic"
)

const needCleanupConcurrencyCount = -1

type Pool struct {
	maxConcurrencyCount int32
	wg                  WaitGroup
	initOnce            sync.Once
	taskChan            chan func()
	manageChan          chan struct{}
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
	case p.manageChan <- struct{}{}:
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
	p.taskChan <- f
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Cleanup() {
	p.init()
	p.wg.Add(1)
	close(p.taskChan)
	p.setMaxConcurrencyCount(needCleanupConcurrencyCount)

	p.wg.Wait()
}

func (p *Pool) init() {
	p.initOnce.Do(func() {
		p.taskChan = make(chan func())
		p.manageChan = make(chan struct{}, 1)
		// pump mange routine working
		p.manageChan <- struct{}{}
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
			case <-p.manageChan:
				break Exit
			case f, ok := <-p.taskChan:
				if !ok {
					break Exit
				}
				p.dispatchWork(f)
			}
		}
	}
	p.wg.Done()
}

func (p *Pool) parkHere() bool {
	for {
		cnt := atomic.LoadInt32(&p.maxConcurrencyCount)
		if cnt == needCleanupConcurrencyCount {
			return true
		}
		if cnt == 0 {
			<-p.manageChan
			continue
		}
		break
	}
	return false
}

func (p *Pool) dispatchWork(f func()) {
	// plus one because manage routine can execute work too.
	if p.wg.CurrentConcurrency()+1 < p.GetMaxConcurrencyCount() {
		p.wg.Go(func() { p.workRoutine(f) })
	} else {
		p.wg.Add(1)
		f()
		p.wg.Done()
	}
}

func (p *Pool) workRoutine(f func()) {
	f()
Exit:
	for {
		if p.WorkerNeedExit() {
			break
		}
		select {
		case f, ok := <-p.taskChan:
			if !ok {
				break Exit
			}
			f()
		default:
			break Exit
		}
	}
}

func (p *Pool) WorkerNeedExit() bool {
	return p.GetCurrentConcurrencyCount() >= p.GetMaxConcurrencyCount()
}
