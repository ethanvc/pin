package pin

import "net/http"

type RouteGroup struct {
	parentGroup     *RouteGroup
	method          string
	relativePath    string
	childGroups     []*RouteGroup
	interceptorFunc []InterceptorFunc
	handler         any
}

func (this *RouteGroup) Group(relativePath string, interceptorFunc ...InterceptorFunc) *RouteGroup {
	child := &RouteGroup{
		parentGroup:     this,
		relativePath:    relativePath,
		interceptorFunc: interceptorFunc,
	}
	this.childGroups = append(this.childGroups, child)
	return child
}

func (this *RouteGroup) Get(relativePath string, handler any, interceptorFunc ...InterceptorFunc) {
	this.handle(http.MethodGet, relativePath, handler, interceptorFunc...)
}

func (this *RouteGroup) Post(relativePath string, handler any, interceptorFunc ...InterceptorFunc) {
	this.handle(http.MethodGet, relativePath, handler, interceptorFunc...)
}

func (this *RouteGroup) handle(method string, relativePath string, handler any, interceptorFunc ...InterceptorFunc) *RouteGroup {
	child := &RouteGroup{
		parentGroup:     this,
		method:          method,
		relativePath:    relativePath,
		interceptorFunc: interceptorFunc,
		handler:         handler,
	}
	this.childGroups = append(this.childGroups, child)
	return child
}
