package pin

import "github.com/ethanvc/pin/pin/status"

type Router struct {
	routeNode routeNode
}

func (this *Router) AddRoute(method string, urlPath string, handler any, interceptorFunc []InterceptorFunc) *status.Status {
	return nil
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

func (this *routeNode) add(method string, urlPath string, handler any, interceptorFunc []InterceptorFunc) *status.Status {
	if len(this.children) == 0 {
		this.children = append(this.children, routeNode{
			commonPath:      urlPath,
			method:          method,
			interceptorFunc: interceptorFunc,
		})
		return nil
	}
	return nil
}

func (this *routeNode) addChild(method string, urlPath string, handler any, interceptorFunc []InterceptorFunc) *status.Status {
	return nil
}

func (this *routeNode) Find(method string, urlPath string, params *Params) (*routeNode, *status.Status) {
	return nil, nil
}
