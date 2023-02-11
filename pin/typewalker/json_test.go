package typewalker

import (
	"github.com/ethanvc/pin/pin/base"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyStruct(t *testing.T) {
	assert.Equal(t, base.StructToJsonStr(struct{}{}), ToLogJsonStr(struct{}{}))
}

func TestFeatureStruct(t *testing.T) {
	type Embed struct {
		E1 int
	}
	type Embed1 struct {
		E2 int
	}
	type TestS struct {
		X  int
		X1 string
		X2 []byte
		X3 []int
		X4 int `json:"xx4"`
		Embed
	}
	v := TestS{
		X:  3,
		X1: "hello world",
		X2: []byte("hello world2"),
		X3: []int{3, 4, 5},
	}
	assert.Equal(t, base.StructToJsonStr(v), ToLogJsonStr(v))
}

func TestEmbeddedStructWithTagName(t *testing.T) {
	type Embed struct {
		E1 int
	}
	type S1 struct {
		Embed `json:"embed"`
	}
	v := S1{}
	assert.Equal(t, base.StructToJsonStr(v), ToLogJsonStr(v))
}

func TestSliceType(t *testing.T) {
	type S struct {
		X []int
	}
	v := S{X: []int{1, 2}}
	assert.Equal(t, base.StructToJsonStr(v), ToLogJsonStr(v))
}

func TestArrayType(t *testing.T) {
	type S struct {
		X [2]int
		Y [2]byte
	}
	v := S{X: [...]int{8, 0}}
	assert.Equal(t, base.StructToJsonStr(v), ToLogJsonStr(v))
}

func TestMapType(t *testing.T) {
	v := map[string]int{
		"hello": 3,
	}
	assert.Equal(t, base.StructToJsonStr(v), ToLogJsonStr(v))
}

func TestMapTYpeDigitKey(t *testing.T) {
	v := map[int]int{
		4: 3,
	}
	assert.Equal(t, base.StructToJsonStr(v), ToLogJsonStr(v))
}
