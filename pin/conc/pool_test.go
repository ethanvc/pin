package conc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
			fmt.Println(token)
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
