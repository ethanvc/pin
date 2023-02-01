package conc

import (
	"sync"
	"sync/atomic"
)

type Pool struct {
	maxConcurrencyCount     int32
	currentConcurrencyCount int32
	nullTaskSet             bool
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

func (p *Pool) Go(f func()) {
	p.init()
	p.task <- f
}

func (p *Pool) Cleanup() {
	p.init()
	close(p.task)
	close(p.manageTask)
	p.waitG.Wait()
}

func (p *Pool) init() {
	p.initGuard.Do(func() {
		p.task = make(chan func())
		p.manageTask = make(chan struct{})
		// block Go until mange routine initialized.
		p.fillNullTask()
		// pump mange routine working
		p.manageTask <- struct{}{}
		go p.manageRoutine()
	})
}

func (p *Pool) manageRoutine() {
	for {
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
			p.fillNullTask()
			<-p.manageTask
			continue
		}
		break
	}
	p.removeNullTask()
}

func (p *Pool) fillNullTask() {
	if p.nullTaskSet {
		return
	}
	p.task <- func() {}
}

func (p *Pool) removeNullTask() {
	if p.nullTaskSet {
		<-p.task
		p.nullTaskSet = false
	}
}

func (p *Pool) dispatchWork(f func()) {
	p.waitG.Add(1)
	if atomic.LoadInt32(&p.currentConcurrencyCount)+1 < atomic.LoadInt32(&p.maxConcurrencyCount) {
		go p.workRoutine(f)
	} else {
		p.workRoutine(f)
		p.waitG.Done()
	}
}

func (p *Pool) workRoutine(f func()) {
	defer p.waitG.Done()
	f()
	for {
		select {
		case f = <-p.task:
			f()
		default:
			break
		}
	}
}
