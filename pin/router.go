package pin

import (
	"github.com/ethanvc/pin/pin/status"
	"github.com/ethanvc/pin/pin/status/codes"
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
		status := this.routeNode.add(realPath, methodHandler{
			method:          method,
			PatternPath:     patternPath,
			handler:         h,
			interceptorFunc: interceptorFunc,
		})
		if status.NotOk() {
			return status
		}
	}
	return nil
}

type methodHandler struct {
	method          string
	PatternPath     string
	interceptorFunc []InterceptorFunc
	handler         Handler
}

type routeNode struct {
	part     string
	children []routeNode

	methodHandlers []methodHandler
}

func findCommLen(s1, s2 string) int {
	l := len(s1)
	if len(s2) < l {
		l = len(s2)
	}
	for i := 0; i < l; i++ {
		if s1[i] == s2[i] {
			continue
		} else {
			return i
		}
	}
	return l
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

func (this *routeNode) add(part string, handler methodHandler) *status.Status {
	if len(this.part) == 0 {
		this.part = part
		this.methodHandlers = []methodHandler{handler}
		return this.initNewNode()
	} else {
		return nil
	}
}

func (this *routeNode) initNewNode() *status.Status {
	if this.part[0] == ':' {
		slashIndex := strings.IndexByte(this.part, '/')
		if slashIndex == -1 {
			return nil
		}
		return this.splitPart(slashIndex)
	}
	if this.part[0] == '*' {
		slashIndex := strings.IndexByte(this.part, '/')
		if slashIndex != -1 {
			return status.NewStatus(codes.Internal, "StarWildcardInvalidHere")
		}
		return this.splitPart(slashIndex)
	}

	partLen := len(this.part)
	for i := 0; i < partLen; i++ {
		if this.part[i] == '/' && i+1 < partLen {
			if this.part[i+1] == ':' || this.part[i+1] == '*' {
				return this.splitPart(i + 1)
			}
		}
	}
	return nil

}

func (this *routeNode) splitPart(sep int) *status.Status {
	child := *this
	*this = routeNode{}
	this.part = child.part[0:sep]
	child.part = child.part[sep:]
	this.children = append(this.children, child)
	return this.children[0].initNewNode()
}
