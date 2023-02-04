package typewalker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStruct struct {
}

func TestToJsonStr(t *testing.T) {
	var s testStruct
	result := ""
	assert.Equal(t, result, ToJsonStr(s))
}
