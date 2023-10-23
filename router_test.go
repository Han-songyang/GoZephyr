package GoZephyr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRouter_AddRoute(t *testing.T) {
	testRouters := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/path1/path2",
		},
	}

	var mockHandler HandleFunc = func(ctx *Context) {

	}
	r := newRouter()
	for _, route := range testRouters {
		r.AddRoute(route.method, route.path, mockHandler)
	}

	// check
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"path1": {
						path: "path1",
						children: map[string]*node{
							"path2": {
								path:     "path2",
								children: map[string]*node{},
								handler:  mockHandler,
							},
						},
					},
				},
			},
		},
	}
	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)
}

func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		n, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的路由"), false
		}
		msg, equal := v.equal(n)
		if !equal {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点数量不匹配"), false
	}
	for path, c := range n.children {
		n, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点 %s 不存在", n.path), false
		}
		msg, ok := c.equal(n)
		if !ok {
			return msg, false
		}
	}
	return "", true
}
