package pin

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRouteGroup_BuildRouter(t *testing.T) {
	controller := &testHandlerStruct{}
	var g RouteGroup
	g.GET("/a/b/c", controller)
	r, status := g.BuildRouter()
	assert.Nil(t, status)
	assert.Equal(t, "/a/b/c/create", r.routeNode.part)
	n := r.Find(http.MethodGet, "/a/b/c/create", nil)
	assert.True(t, n.ValidHandler())
}
