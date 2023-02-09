package typewalker

import (
	"github.com/ethanvc/pin/pin/base"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNil(t *testing.T) {
	assert.Equal(t, base.StructToJsonStr(nil), ToLogJsonStr(nil))
}

func TestEmptyStruct(t *testing.T) {
	assert.Equal(t, base.StructToJsonStr(struct{}{}), ToLogJsonStr(struct{}{}))
}

func TestNonEmptyStruct(t *testing.T) {
	type TestS struct {
		X  int
		X1 int8
		X2 int16
		X3 int32
		X4 int64
		Y  uint
		Y1 uint8
		Y2 uint16
		Y3 uint32
		Y4 uint64
		Y5 byte
	}
	v := TestS{
		X: 3,
	}
	assert.Equal(t, base.StructToJsonStr(v), ToLogJsonStr(v))
}
