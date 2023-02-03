package base

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIn(t *testing.T) {
	assert.Equal(t, true, In(0, 0, 1, 2))
	assert.Equal(t, false, In(0, 1, 2))
}
