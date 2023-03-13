package pin

import (
	"github.com/ethanvc/pin/pin/base"
	"github.com/ethanvc/pin/pin/status"
	"github.com/ethanvc/pin/pin/status/codes"
	"github.com/iancoleman/strcase"
	"path"
	"strings"
)

type Router struct {
	routeNode routeNode
}

func (this *Router) addRoute(method string, patternPath string, handler any,
	interceptorFunc []InterceptorFunc) *status.Status {
	handlers := NewHandlers(handler)
	for _, h := range handlers {
		realPath := patternPath
		if len(h.methodName) > 0 {
			realPath = path.Join(realPath, strcase.ToKebab(h.methodName))
		}
		status := this.routeNode.add(realPath, HttpHandler{
			Method:          method,
			PatternPath:     patternPath,
			Handler:         h,
			InterceptorFunc: interceptorFunc,
		})
		if status.NotOk() {
			return status
		}
	}
	return nil
}

func (this *Router) Find(method, urlPath string, params *Params) HttpHandler {
	*params = (*params)[:0]
	return this.routeNode.Find(method, urlPath, params)
}

type HttpHandler struct {
	Method          string
	PatternPath     string
	InterceptorFunc []InterceptorFunc
	Handler         Handler
}

func (this HttpHandler) IsValid() bool {
	return this.Handler.IsValid()
}

type httpHandlers []HttpHandler

func (this httpHandlers) Find(method string) (HttpHandler, bool) {
	for _, h := range this {
		if h.Method == HttpMethodAny || h.Method == method {
			return h, true
		}
	}
	return HttpHandler{}, false
}

type routeNode struct {
	part     string
	children []routeNode

	httpHandlers httpHandlers
}

func commLen(s1, s2 string) int {
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

func (this *routeNode) add(part string, handler HttpHandler) *status.Status {
	if len(this.part) == 0 {
		this.part = part
		this.httpHandlers = []HttpHandler{handler}
		return this.initNewNode()
	}

	commL := commLen(part, this.part)
	part = part[commL:]
	if len(part) == 0 {
		if _, ok := this.httpHandlers.Find(handler.Method); ok {
			return status.NewStatus(codes.Internal, "RouteConflict")
		}
		this.httpHandlers = append(this.httpHandlers, handler)
		return nil
	}

	if this.wildcardNode() && part[0] != '/' {
		return status.NewStatus(codes.Internal, "WildcardDuplicate")
	}

	if commL == len(this.part) {
		for i := 0; i < len(this.children); i++ {
			commL := commLen(this.children[i].part, part)
			if commL > 0 {
				return this.children[i].add(part, handler)
			}
		}
		this.children = append(this.children, routeNode{
			part:         part,
			httpHandlers: []HttpHandler{handler},
		})
		return this.children[len(this.children)-1].initNewNode()
	}

	newRoot := routeNode{
		part: this.part[0:commL],
	}
	oldChild := *this
	oldChild.part = this.part[commL:]
	newChild := routeNode{
		part:         part,
		httpHandlers: httpHandlers{handler},
	}
	status := newChild.initNewNode()
	if status.NotOk() {
		return status
	}
	newRoot.children = append(newRoot.children, oldChild, newChild)
	*this = newRoot
	return nil
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

func (this *routeNode) Find(method, part string, params *Params) HttpHandler {
	if this.part[0] == '*' {
		if len(this.httpHandlers) == 0 {
			return HttpHandler{}
		}
		if h, ok := this.httpHandlers.Find(method); ok {
			*params = append(*params, Param{
				Key:   this.part[1:],
				Value: part,
			})
			return h
		}
		return HttpHandler{}
	}
	if this.part[0] == ':' {
		slashIndex := strings.IndexByte(part, '/')
		if slashIndex == -1 {
			if h := this.findHandler(method); h.IsValid() {
				*params = append(*params, Param{
					Key:   this.part[1:],
					Value: part,
				})
				return h
			}
			return HttpHandler{}
		}
		*params = append(*params, Param{
			Key:   this.part[1:],
			Value: part[0:slashIndex],
		})
		part = part[slashIndex:]
		return this.findInChildren(method, part, params)
	}

	if !strings.HasPrefix(part, this.part) {
		return HttpHandler{}
	}
	part = part[len(this.part):]
	if len(part) == 0 {
		return this.findHandler(method)
	}
	return this.findInChildren(method, part, params)
}

func (this *routeNode) findInChildren(method, part string, params *Params) HttpHandler {
	for i := 0; i < len(this.children); i++ {
		if this.children[i].canStepIn(part) {
			return this.children[i].Find(method, part, params)
		}
	}
	return HttpHandler{}
}

func (this *routeNode) findHandler(method string) HttpHandler {
	if h, ok := this.httpHandlers.Find(method); ok {
		return h
	} else {
		return HttpHandler{}
	}
}

func (this *routeNode) wildcardNode() bool {
	return base.In(this.part[0], '*', ':')
}

func (this *routeNode) canStepIn(part string) bool {
	if this.wildcardNode() || strings.HasPrefix(part, this.part) {
		return true
	} else {
		return false
	}
}
