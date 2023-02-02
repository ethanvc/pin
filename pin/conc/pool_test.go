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
	verify := func(maxConcurrencyCount, currentCurrencyCount int) {
		assert.Equal(t, maxConcurrencyCount, pool.GetMaxConcurrencyCount())
		assert.Equal(t, currentCurrencyCount, pool.GetCurrentConcurrencyCount())
	}
	verify(0, 0)
	pool.SetMaxConcurrencyCount(1)
	verify(1, 0)
	wg := runTask(1000, 1, 0, &pool)
	wg.Wait()
	verify(1, 0)
}
