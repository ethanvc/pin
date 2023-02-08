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
		X int
	}
	v := TestS{
		X: 3,
	}
	assert.Equal(t, base.StructToJsonStr(v), ToLogJsonStr(v))
}
