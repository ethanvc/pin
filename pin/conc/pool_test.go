package conc

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"math"
	"sync/atomic"
	"testing"
	"time"
)

type RunTask struct {
	TaskCount           int64
	SendTaskConcurrency int
	Pool                *Pool
	Interval            time.Duration
	WaitGroup           WaitGroup
	T                   *testing.T
}

func (r *RunTask) Verify(maxConcurrencyCount, currentCurrencyCount int) {
	assert.Equal(r.T, maxConcurrencyCount, r.Pool.GetMaxConcurrencyCount())
	assert.Equal(r.T, currentCurrencyCount, r.Pool.GetCurrentConcurrencyCount())
	time.Sleep(1 * time.Millisecond)
	assert.Equal(r.T, maxConcurrencyCount, r.Pool.GetMaxConcurrencyCount())
	assert.Equal(r.T, currentCurrencyCount, r.Pool.GetCurrentConcurrencyCount())
}

func (r *RunTask) Start() {
	for i := 0; i < r.SendTaskConcurrency; i++ {
		r.WaitGroup.Go(func() {
			r.worker()
		})
	}
}

func (r *RunTask) worker() {
	for {
		token := atomic.AddInt64(&r.TaskCount, -1)
		if token < 0 {
			break
		}
		r.Pool.Go(func() {
			time.Sleep(r.Interval)
			// fmt.Println(token)
		})
	}
}

func Prepare(t *testing.T) (*Pool, *RunTask) {
	var pool Pool
	runner := RunTask{
		TaskCount:           1000,
		SendTaskConcurrency: 10,
		Pool:                &pool,
		Interval:            0,
		T:                   t,
	}
	return &pool, &runner
}

func TestPool(t *testing.T) {
	defer goleak.VerifyNone(t)
	p, r := Prepare(t)
	r.Start()
	r.Verify(0, 0)
	assert.Equal(t, int64(990), r.TaskCount)
	p.SetMaxConcurrencyCount(3)
	r.WaitGroup.Wait()
	p.Wait()
	p.Cleanup()
	r.Verify(needCleanupConcurrencyCount, 0)
}

// will crash.
func testCleanupWhenRun(t *testing.T) {
	defer goleak.VerifyNone(t)
	p, r := Prepare(t)
	r.TaskCount = math.MaxInt64
	r.Start()
	time.Sleep(time.Millisecond)
	p.Cleanup()
	r.TaskCount = 0
	r.WaitGroup.Wait()
	r.Verify(needCleanupConcurrencyCount, 0)
}

func TestConcurrencyReduce(t *testing.T) {
	defer goleak.VerifyNone(t)
	p, r := Prepare(t)
	r.TaskCount = math.MaxInt64
	r.Interval = time.Millisecond
	r.Start()
	time.Sleep(time.Millisecond)
	p.SetMaxConcurrencyCount(50)
	for {
		if p.GetCurrentConcurrencyCount() == 50 {
			break
		}
		time.Sleep(0)
	}
	p.SetMaxConcurrencyCount(10)
	for {
		if p.GetCurrentConcurrencyCount() == 10 {
			break
		}
		time.Sleep(0)
	}
	p.SetMaxConcurrencyCount(0)
	for {
		if p.GetCurrentConcurrencyCount() == 0 {
			break
		}
		time.Sleep(0)
	}
	p.SetMaxConcurrencyCount(100)
	r.TaskCount = 0
	r.WaitGroup.Wait()
	p.Cleanup()
}

func TestConcurrencyIncrease(t *testing.T) {
	defer goleak.VerifyNone(t)
	p, r := Prepare(t)
	r.Interval = time.Millisecond
	r.TaskCount = math.MaxInt64
	r.Start()
	time.Sleep(time.Millisecond)
	p.SetMaxConcurrencyCount(50)
	for {
		if p.GetCurrentConcurrencyCount() == 50 {
			break
		}
		time.Sleep(0)
	}
	p.SetMaxConcurrencyCount(60)
	for {
		if p.GetCurrentConcurrencyCount() == 60 {
			break
		}
		time.Sleep(0)
	}
	r.TaskCount = 0
	r.WaitGroup.Wait()
	p.Cleanup()
}
