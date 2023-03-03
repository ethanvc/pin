package pin

import (
	"context"
	"github.com/ethanvc/pin/pin/status"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRouteGroup_BuildRouter(t *testing.T) {
	emptyFunc := func(context.Context, *int) (*int, *status.Status) { return nil, nil }
	var g RouteGroup
	g.GET("/a/:b/c", emptyFunc)
	g.GET("/a/:b/c/d", emptyFunc)
	r, _ := g.BuildRouter()
	var params Params
	h := r.Find(http.MethodGet, "/a/c/c/d", &params)
	assert.Equal(t, "/a/:b/c/d", h.PatternPath)
}
