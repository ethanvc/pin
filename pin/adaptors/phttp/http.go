package phttp

import (
	"github.com/ethanvc/pin/pin"
	"github.com/ethanvc/pin/pin/status"
	"golang.org/x/net/context"
	"net/http"
)

func ServeHttp(svr *pin.Server, route *pin.Router, w http.ResponseWriter, r *http.Request) {
	pr := HttpProtocolRequest{
		route: route,
		w:     w,
		r:     r,
	}
	svr.ProcessRequest(pr)
}

type HttpProtocolRequest struct {
	route *pin.Router
	w     http.ResponseWriter
	r     *http.Request
}

func (this HttpProtocolRequest) InitializeRequest(req *pin.Request) *status.Status {
	var params *pin.Params
	h := this.route.Find(this.r.Method, this.r.URL.Path, params)
	req.PatternPath = this.r.Method + ":" + h.PatternPath
	req.Handler = h.Handler
	req.Interceptors = h.InterceptorFunc
	return nil
}

func (this HttpProtocolRequest) FinalizeRequest(req *pin.Request) *status.Status {
	return nil
}

func (this HttpProtocolRequest) ProtocolContext() context.Context {
	return this.r.Context()
}

func (this HttpProtocolRequest) ProtocolObject() any {
	return this
}
