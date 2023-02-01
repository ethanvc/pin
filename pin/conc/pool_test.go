package conc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPool(t *testing.T) {
	var pool Pool
	assert.Equal(t, 0, pool.GetMaxConcurrencyCount())
	assert.Equal(t, 0, pool.GetCurrentConcurrencyCount())
}
