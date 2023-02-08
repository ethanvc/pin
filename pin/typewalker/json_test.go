package typewalker

import (
	"github.com/ethanvc/pin/pin/base"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStruct struct {
	F  int
	F2 int `json:"f2"`
}

func TestNil(t *testing.T) {
	assert.Equal(t, base.StructToJsonStr(nil), ToLogJsonStr(nil))
}

func TestStruct(t *testing.T) {
	var s testStruct
	s.F = 3
	result := ""
	assert.Equal(t, result, ToLogJsonStr(s))
}
