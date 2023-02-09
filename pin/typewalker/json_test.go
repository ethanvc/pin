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

func TestFeatureStruct(t *testing.T) {
	type TestS struct {
		X  int
		X1 string
	}
	v := TestS{
		X:  3,
		X1: "hello world",
	}
	assert.Equal(t, base.StructToJsonStr(v), ToLogJsonStr(v))
}
