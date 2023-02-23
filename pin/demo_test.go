package pin

import (
	"context"
	"github.com/ethanvc/pin/pin/status"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestController struct {
}

type TestReq struct {
	Name string
}

type TestResp struct {
	Name string
}

func (this *TestController) Get(c context.Context, req *TestReq) (*TestResp, *status.Status) {
	return &TestResp{
		Name: req.Name,
	}, nil
}

func TestDemo(t *testing.T) {
	var controller TestController
	req := &TestReq{
		Name: "hello",
	}
	resp, status := CreatePlainCall("ControllerGet", controller.Get).Call(context.Background(), req)
	assert.Nil(t, status)
	assert.Equal(t, req.Name, resp.Name)
}
