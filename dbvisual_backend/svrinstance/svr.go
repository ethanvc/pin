package svrinstance

import (
	"github.com/ethanvc/pin/dbvisual_backend/controller"
	"github.com/ethanvc/pin/pin"
	"github.com/ethanvc/pin/pin/adaptors/phttp"
	"github.com/ethanvc/pin/pin/plog"
	"net/http"
)

func Serve() {
	g := pin.RouteGroup{}
	api := g.Group("/db-visual/api")
	api.POST("/instance", &controller.Instance{})
	route := api.MustBuildRouter()

	s := &pin.Server{}
	handler := func(w http.ResponseWriter, r *http.Request) {
		phttp.ServeHttp(s, route, w, r)
	}
	addr := "127.0.0.1:8011"
	plog.Default().Info("StartSvr").Str("svr_addr", addr).Done()
	http.ListenAndServe(addr, http.HandlerFunc(handler))
}
