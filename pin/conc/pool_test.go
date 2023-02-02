package conc

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func runTask(taskCount int, sendTaskConcurrency int, sleepInterval time.Duration, p *Pool) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(taskCount)
	taskCounter := int32(0)
	for i := 0; i < sendTaskConcurrency; i++ {
		go func() {
			for {
				currentCounter := atomic.AddInt32(&taskCounter, 1)
				if int(currentCounter) > taskCount {
					break
				}
				p.Go(func() {
					time.Sleep(sleepInterval)
					wg.Done()
				})
			}
		}()
	}
	return &wg
}

func TestPool(t *testing.T) {
	var pool Pool
	defer pool.Cleanup()
	assert.Equal(t, 0, pool.GetMaxConcurrencyCount())
	assert.Equal(t, 0, pool.GetCurrentConcurrencyCount())
	pool.SetMaxConcurrencyCount(1)
	assert.Equal(t, 1, pool.GetMaxConcurrencyCount())
	wg := runTask(1000, 1, 0, &pool)
	wg.Wait()
	assert.Equal(t, 0, pool.GetCurrentConcurrencyCount())
}
