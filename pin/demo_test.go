package pin

import (
	"context"
	"github.com/ethanvc/pin/pin/status"
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
	var s Server
	req := CreatePlainCall("Get", func() *status.Status {
		return nil
	})
	s.ProcessRequest(req)
}
