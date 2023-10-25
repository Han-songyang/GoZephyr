package GoZephyr

import (
	"fmt"
	"strings"
)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
type Params []Param

func (ps Params) Get(name string) (string, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			return entry.Value, true
		}
	}
	return "", false
}

// router routing tree
type router struct {
	// map[method]*node
	// key:   HTTP method,One http method corresponds to one tree
	// value: root node
	trees map[string]*node
}

// AddRoute is a method to Create nodes, including paths and HandleFunc
// and register it in the routing tree.
func (r *router) addRoute(method string, path string, handler HandleFunc) {
	if path == "" {
		panic("web: Route is an empty string")
	}
	if path[0] != '/' {
		panic("web: Routes must start with '/'")
	}

	if path != "/" && path[len(path)-1] == '/' {
		panic("web: Routes cannot end in '/'")
	}

	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path:     "/",
			children: map[string]*node{},
		}
		r.trees[method] = root
	}
	if path == "/" {
		if root.handler != nil {
			panic(fmt.Sprintf("web: Route has been registered [%s]", path))
		}
		root.handler = handler
		return
	}
	root.addPath(path, handler)
}

// findRoute is a method to find the node in the routing tree
func (r *router) findRoute(method string, path string) (*nodeInfo, bool) {
	ni := new(nodeInfo)
	var nType nodeType
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		ni.n = root
		return ni, true
	}

	paths := strings.Split(strings.Trim(path, "/"), "/")
	for _, s := range paths {
		root, nType, ok = root.childOf(s)
		if !ok {
			return nil, false
		}
		if nType == param {
			ni.addValue(root.path[1:], s)
		}
	}

	ni.n = root
	return ni, true
}

type nodeType uint8

const (
	static nodeType = iota + 1
	wildcard
	param
)

// node routing tree nodes
type node struct {
	//
	nodeType nodeType
	// routing path
	path string

	// children child node
	// map[path]*node
	children map[string]*node

	// wildcardChild wildcard node: /path/*
	wildcardChild *node

	// param param node: /path/:name
	paramChild *node

	// handler Method of registering after hitting a path
	handler HandleFunc
}

// addPath is a method to Create a node, including paths and HandleFunc
func (n *node) addPath(path string, handler HandleFunc) {
	paths := strings.Split(path[1:], "/")
	root := n
	for _, p := range paths {
		if p == "" {
			panic(fmt.Sprintf("path error [%s]", path))
		}
		root = root.createPathNode(p)
	}
	if root.handler != nil {
		panic(fmt.Sprintf("web: Route has been registered [%s]", path))
	}
	root.handler = handler
}

// createPathNode is a method to Create Child node, and return child node
func (n *node) createPathNode(path string) *node {
	// wildcard node
	if path == "*" {
		if n.wildcardChild == nil {
			n.wildcardChild = &node{
				path:     path,
				nodeType: wildcard,
			}
		}
		return n.wildcardChild
	}

	// param node
	if path[0] == ':' {
		if n.paramChild == nil {
			n.paramChild = &node{
				path:     path,
				nodeType: param,
			}
		} else if n.paramChild.path != path {
			panic(fmt.Sprintf("web: Route has been registered [%s]", path))
		}
		return n.paramChild
	}

	// static route
	if n.children == nil {
		n.children = make(map[string]*node)
	}

	child, ok := n.children[path]
	if !ok {
		child = &node{
			path:     path,
			nodeType: static,
		}
		n.children[path] = child
	}
	return child
}

func (n *node) childOf(path string) (*node, nodeType, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, param, true
		}
		return n.wildcardChild, wildcard, n.wildcardChild != nil
	}
	res, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, param, true
		}
		return n.wildcardChild, wildcard, n.wildcardChild != nil
	}
	return res, static, ok
}

type nodeInfo struct {
	n      *node
	params Params
}

func (m *nodeInfo) addValue(key string, value string) {
	if m.params == nil {
		// 大多数情况，参数路径只会有一段
		m.params = Params{
			Param{Key: key, Value: value},
		}
		return
	}
	m.params = append(m.params, Param{
		Key:   key,
		Value: value,
	})
	return
}

// newRouter return a new route
func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}
