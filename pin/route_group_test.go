package pin

import (
	"context"
	"github.com/ethanvc/pin/pin/status"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouteGroup_BuildRouter(t *testing.T) {
	emptyFunc := func(context.Context, *int) (*int, *status.Status) { return nil, nil }
	var g RouteGroup
	g.GET("/a/:b/c", emptyFunc)
	r, status := g.BuildRouter()
	_ = r
	assert.Nil(t, status)
}
