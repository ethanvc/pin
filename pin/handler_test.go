package pin

import (
	"context"
	"github.com/ethanvc/pin/pin/status"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testHandlerStruct struct {
}

func (*testHandlerStruct) Get(c context.Context, req *testHandlerStruct) (*testHandlerStruct, *status.Status) {
	return nil, nil
}

func (*testHandlerStruct) Create(c context.Context, req *testHandlerStruct) (*testHandlerStruct, *status.Status) {
	return nil, nil
}

func (*testHandlerStruct) get(c context.Context, req *testHandlerStruct) (*testHandlerStruct, *status.Status) {
	return nil, nil
}

func (*testHandlerStruct) List(c context.Context, req testHandlerStruct) (*testHandlerStruct, *status.Status) {
	return nil, nil
}

func TestNewHandlers(t *testing.T) {
	v := &testHandlerStruct{}
	handlers := NewHandlers(v)
	assert.Equal(t, 2, len(handlers))
}
