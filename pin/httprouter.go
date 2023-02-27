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

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

func (this *HttpRouter) add(method string, urlPath string, handler any, interceptorFunc []InterceptorFunc) *status.Status {
	return nil
}

func (this *HttpRouter) Find(method string, urlPath string, params *Params) (*routeNode, *status.Status) {
	return nil, nil
}
