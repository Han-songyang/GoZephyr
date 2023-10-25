package GoZephyr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
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
			path:   "/path",
		},
		{
			method: http.MethodGet,
			path:   "/path/path2",
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
		// 通配符测试用例
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc/*",
		},
		// 参数路由
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/detail",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/*",
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
					"param": {
						path: "param",
						paramChild: &node{
							path: ":id",
							wildcardChild: &node{
								path:    "*",
								handler: mockHandler,
							},
							children: map[string]*node{"detail": {path: "detail", handler: mockHandler}},
							handler:  mockHandler,
						},
					},
					"path": {
						path: "path",
						children: map[string]*node{
							"path2": {
								path:    "path2",
								handler: mockHandler,
							},
						},
						handler: mockHandler,
					},
					"order": {
						path: "order",
						children: map[string]*node{
							"detail": {
								path:    "detail",
								handler: mockHandler,
							},
						},
						wildcardChild: &node{
							path:     "*",
							handler:  mockHandler,
							children: map[string]*node{},
						},
					},
				},
				wildcardChild: &node{
					path: "*",
					children: map[string]*node{
						"abc": {
							path: "abc",
							wildcardChild: &node{
								path: "*", handler: mockHandler,
								children: map[string]*node{},
							},
							handler: mockHandler,
						},
					},
					wildcardChild: &node{
						path:     "*",
						handler:  mockHandler,
						children: map[string]*node{},
					},
					handler: mockHandler,
				},
				handler: mockHandler,
			},
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
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
						path:    "login",
						handler: mockHandler,
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
	assert.PanicsWithValue(t, "web: Route has been registered [:id]", func() {
		r.addRoute(http.MethodGet, "/param/:name", mockHandler)
		r.addRoute(http.MethodGet, "/param/:id", mockHandler)
	})
}

func Test_router_findRoute(t *testing.T) {
	testRoutes := []struct {
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
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodGet,
			path:   "/user/*/home",
		},
		{
			method: http.MethodPost,
			path:   "/order/*",
		},
		// 参数路由
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/detail",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/*",
		},
		{
			method: http.MethodGet,
			path:   "/path/:id/:name",
		},
		{
			method: http.MethodGet,
			path:   "/:id",
		},
	}

	mockHandler := func(ctx *Context) {}

	testCases := []struct {
		name   string
		method string
		path   string
		found  bool
		mi     *nodeInfo
	}{
		{
			name:   "method not found",
			method: http.MethodHead,
		},
		//{
		//	name:   "path not found",
		//	method: http.MethodGet,
		//	path:   "/abc",
		//},
		{
			name:   "root",
			method: http.MethodGet,
			path:   "/",
			found:  true,
			mi: &nodeInfo{
				n: &node{
					path:    "/",
					handler: mockHandler,
				},
			},
		},
		{
			name:   "user",
			method: http.MethodGet,
			path:   "/user",
			found:  true,
			mi: &nodeInfo{
				n: &node{
					path:    "user",
					handler: mockHandler,
				},
			},
		},
		{
			name:   "no handler",
			method: http.MethodPost,
			path:   "/order",
			found:  true,
			mi: &nodeInfo{
				n: &node{
					path: "order",
				},
			},
		},
		{
			name:   "two layer",
			method: http.MethodPost,
			path:   "/order/create",
			found:  true,
			mi: &nodeInfo{
				n: &node{
					path:    "create",
					handler: mockHandler,
				},
			},
		},
		// 通配符匹配
		{
			// 命中/order/*
			name:   "star match",
			method: http.MethodPost,
			path:   "/order/delete",
			found:  true,
			mi: &nodeInfo{
				n: &node{
					path:    "*",
					handler: mockHandler,
				},
			},
		},
		{
			// 命中通配符在中间的
			// /user/*/home
			name:   "star in middle",
			method: http.MethodGet,
			path:   "/user/Tom/home",
			found:  true,
			mi: &nodeInfo{
				n: &node{
					path:    "home",
					handler: mockHandler,
				},
			},
		},
		{
			// 比 /order/* 多了一段
			name:   "overflow",
			method: http.MethodPost,
			path:   "/order/delete/123",
		},
		// 参数匹配
		{
			// 命中 /param/:id
			name:   ":id",
			method: http.MethodGet,
			path:   "/param/123",
			found:  true,
			mi: &nodeInfo{
				n: &node{
					path:    ":id",
					handler: mockHandler,
				},
				params: Params{
					Param{
						Key:   "id",
						Value: "123",
					},
				},
			},
		},
		{
			// 命中 /param/:id/*
			name:   ":id*",
			method: http.MethodGet,
			path:   "/param/123/abc",
			found:  true,
			mi: &nodeInfo{
				n: &node{
					path:    "*",
					handler: mockHandler,
				},
				params: Params{
					Param{
						Key:   "id",
						Value: "123",
					},
				},
			},
		},

		{
			// 命中 /param/:id/detail
			name:   ":id*",
			method: http.MethodGet,
			path:   "/param/123/detail",
			found:  true,
			mi: &nodeInfo{
				n: &node{
					path:    "detail",
					handler: mockHandler,
				},
				params: Params{
					Param{
						Key:   "id",
						Value: "123",
					},
				},
			},
		},
		{
			// 命中 /param/:id/:name
			name:   ":id*",
			method: http.MethodGet,
			path:   "/path/123/han",
			found:  true,
			mi: &nodeInfo{
				n: &node{
					path:    ":name",
					handler: mockHandler,
				},
				params: Params{
					Param{
						Key:   "id",
						Value: "123",
					},
					Param{
						Key:   "name",
						Value: "han",
					},
				},
			},
		},
		{
			// 命中 /:id
			name:   ":id",
			method: http.MethodGet,
			path:   "/123",
			found:  true,
			mi: &nodeInfo{
				n: &node{
					path:    ":id",
					handler: mockHandler,
				},
				params: Params{
					Param{
						Key:   "id",
						Value: "123",
					},
				},
			},
		},
	}

	r := newRouter()
	for _, tr := range testRoutes {
		r.addRoute(tr.method, tr.path, mockHandler)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.found, found)
			if !found {
				return
			}
			assert.Equal(t, tc.mi.params, mi.params)
			n := mi.n
			wantVal := reflect.ValueOf(tc.mi.n.handler)
			nVal := reflect.ValueOf(n.handler)
			assert.Equal(t, wantVal, nVal)
			if mi.n.nodeType == param {
				a, _ := tc.mi.params.Get("id")
				assert.Equal(t, a, "123")
				a, _ = tc.mi.params.Get("aaa")
				assert.Equal(t, a, "")
			}
		})
	}
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
	if y == nil {
		return "目标节点为 nil\n", false
	}
	if n.path != y.path {
		return fmt.Sprintf("节点路径不匹配\n"), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点数量不匹配\n"), false
	}
	if n.wildcardChild != nil {
		str, ok := n.wildcardChild.equal(y.wildcardChild)
		if !ok {
			return fmt.Sprintf("%s 通配符节点不匹配 %s\n", n.path, str), false
		}
	}
	for path, c := range n.children {
		n, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点 %s 不存在\n", n.path), false
		}
		msg, ok := c.equal(n)
		if !ok {
			return msg, false
		}
	}
	return "", true
}
