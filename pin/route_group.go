package pin

import (
	"github.com/ethanvc/pin/pin/status"
	"net/http"
)
import pathf "path"

const HttpMethodAny = "Any"

type RouteGroup struct {
	parentGroup     *RouteGroup
	method          string
	path            string
	children        []*RouteGroup
	interceptorFunc []InterceptorFunc
	handler         any
}

func (this *RouteGroup) Group(path string, interceptorFunc ...InterceptorFunc) *RouteGroup {
	child := &RouteGroup{
		parentGroup:     this,
		path:            this.mergeParentPath(path),
		interceptorFunc: this.mergeParentInterceptors(interceptorFunc),
	}
	this.children = append(this.children, child)
	return child
}

func (this *RouteGroup) mergeParentPath(path string) string {
	if this.parentGroup != nil {
		return pathf.Join(this.parentGroup.path, path)
	} else {
		return path
	}
}

func (this *RouteGroup) mergeParentInterceptors(interceptorFunc []InterceptorFunc) []InterceptorFunc {
	if this.parentGroup != nil {
		tmp := append([]InterceptorFunc{}, this.parentGroup.interceptorFunc...)
		return append(tmp, interceptorFunc...)
	} else {
		return interceptorFunc
	}
}

func (this *RouteGroup) GET(relativePath string, handler any, interceptorFunc ...InterceptorFunc) {
	this.Handle(http.MethodGet, relativePath, handler, interceptorFunc...)
}

func (this *RouteGroup) POST(relativePath string, handler any, interceptorFunc ...InterceptorFunc) {
	this.Handle(http.MethodGet, relativePath, handler, interceptorFunc...)
}

func (this *RouteGroup) Handle(method string, relativePath string, handler any, interceptorFunc ...InterceptorFunc) *RouteGroup {
	child := &RouteGroup{
		parentGroup:     this,
		method:          method,
		path:            relativePath,
		interceptorFunc: interceptorFunc,
		handler:         handler,
	}
	this.children = append(this.children, child)
	return child
}

func (this *RouteGroup) BuildRouter() (*Router, *status.Status) {
	r := &Router{}
	status := this.buildRouter(r)
	if status.NotOk() {
		return r, status
	}
	return r, nil
}

func (this *RouteGroup) buildRouter(r *Router) *status.Status {
	if len(this.method) > 0 {
		status := r.AddRoute(this.method, this.path, this.handler, this.interceptorFunc)
		if status.NotOk() {
			return status
		}
	}

	for _, child := range this.children {
		status := child.buildRouter(r)
		if status.NotOk() {
			return status
		}
	}
	return nil
}
