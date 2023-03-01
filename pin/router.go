package pin

import (
	"github.com/ethanvc/pin/pin/status"
	"github.com/iancoleman/strcase"
	"path"
	"strings"
)

type Router struct {
	routeNode routeNode
}

func (this *Router) AddRoute(method string, patternPath string, handler any,
	interceptorFunc []InterceptorFunc) *status.Status {
	handlers := NewHandlers(handler)
	for _, h := range handlers {
		realPath := patternPath
		if len(h.methodName) > 0 {
			realPath = path.Join(realPath, strcase.ToKebab(h.methodName))
		}
		status := this.routeNode.add(method, realPath, realPath, h, interceptorFunc)
		if status.NotOk() {
			return status
		}
	}
	return nil
}

func (this *Router) Find(method string, urlPath string, params *Params) routeNode {
	return this.routeNode.find(method, urlPath, params)
}

type routeNode struct {
	part            string
	method          string
	children        []routeNode
	interceptorFunc []InterceptorFunc
	handler         Handler
	PatternPath     string
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

func (this *routeNode) add(method string, patternPath string, part string, handler Handler,
	interceptorFunc []InterceptorFunc) *status.Status {
	if len(this.part) == 0 {
		this.method = method
		this.PatternPath = patternPath
		this.part = part
		this.handler = handler
		this.interceptorFunc = interceptorFunc
		return this.makeNewNode()
	}
	return nil
}

func (this *routeNode) addChild(method string, urlPath string, handler any,
	interceptorFunc []InterceptorFunc) *status.Status {
	return nil
}

func (this *routeNode) find(method string, part string, params *Params) routeNode {
	if part == this.part {
		if this.method == method || this.method == HttpMethodAny {
			return *this
		} else {
			return routeNode{}
		}
	} else {
		part = part[len(this.part):]
		if strings.HasPrefix(part, this.part) {
			for _, n := range this.children {
				result := n.find(method, part, params)
				if result.ValidHandler() {
					return result
				}
			}
		}
		return routeNode{}
	}
}

func (this *routeNode) makeNewNode() *status.Status {
	partLen := len(this.part)
	for i := 0; i < partLen; i++ {
		if this.part[i] == '/' && i+1 < partLen {
			if this.part[i+1] == ':' || this.part[i+1] == '*' {
				child := *this
				child.part = child.part[i+1:]
				this.part = this.part[:i+1]
				this.method = ""
				this.interceptorFunc = nil
				this.handler = Handler{}
				this.PatternPath = ""
				this.children = append(this.children, child)
				return this.children[0].makeNewNode()
			}
		}
	}
	return nil
}

func (this *routeNode) ValidHandler() bool {
	return this.handler.methodVal.IsValid()
}
