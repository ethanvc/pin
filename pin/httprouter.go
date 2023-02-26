package pin

import "github.com/ethanvc/pin/pin/status"

type HttpRouter struct {
	children []routeNode
}

type routeNode struct {
	commonPath      string
	method          string
	children        []routeNode
	interceptorFunc []InterceptorFunc
	handler         Handler
}

func CreateHttpRouter(group *RouteGroup) (*HttpRouter, *status.Status) {
	r := &HttpRouter{}
	return r, nil
}
