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
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
	}

	var mockHandler HandleFunc = func(ctx *Context) {

	}
	r := newRouter()
	for _, route := range testRouters {
		r.addRoute(route.method, route.path, mockHandler)
	}

	// check
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"user": {
						path: "user",
						children: map[string]*node{
							"home": {
								path:    "home",
								handler: mockHandler,
							},
						},
						handler: mockHandler},
					"order": {
						path: "order",
						children: map[string]*node{
							"detail": {
								path:    "detail",
								handler: mockHandler,
							},
						},
					},
				},
				handler: mockHandler,
			},
			http.MethodPost: {
				path: "/", children: map[string]*node{
					"order": {
						path: "order",
						children: map[string]*node{
							"create": {
								path:    "create",
								handler: mockHandler,
							},
						},
					},
					"login": {
						path: "login", handler: mockHandler,
					},
				},
			},
		},
	}
	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)

	// 非法用例
	r = newRouter()

	// 空字符串
	assert.PanicsWithValue(t, "web: Route is an empty string", func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	})

	// 前导没有 /
	assert.PanicsWithValue(t, "web: Routes must start with '/'", func() {
		r.addRoute(http.MethodGet, "a/b/c", mockHandler)
	})

	// 后缀有 /
	assert.PanicsWithValue(t, "web: Routes cannot end in '/'", func() {
		r.addRoute(http.MethodGet, "/a/b/c/", mockHandler)
	})

	// 根节点重复注册
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.PanicsWithValue(t, "web: Route has been registered [/]", func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	})
	// 普通节点重复注册
	r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.PanicsWithValue(t, "web: Route has been registered [/a/b/c]", func() {
		r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	})

	// 多个 /
	assert.PanicsWithValue(t, "path error [/a//b]", func() {
		r.addRoute(http.MethodGet, "/a//b", mockHandler)
	})
	assert.PanicsWithValue(t, "path error [//a/b]", func() {
		r.addRoute(http.MethodGet, "//a/b", mockHandler)
	})
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
