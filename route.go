package GoZephyr

import (
	"fmt"
	"strings"
)

// router routing tree
type router struct {
	// map[method]*node
	// key:   HTTP method,One http method corresponds to one tree
	// value: root node
	trees map[string]*node
}

// node routing tree nodes
type node struct {
	// routing path
	path string

	// children child node
	// map[path]*node
	children map[string]*node

	// handler Method of registering after hitting a path
	handler HandleFunc
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
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	if !ok {
		child = &node{path: path}
		n.children[path] = child
	}
	return child
}

func (r *router) findRoute(method string, path string) (*node, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return root, true
	}

	paths := strings.Split(strings.Trim(path, "/"), "/")
	for _, s := range paths {
		root, ok = root.childOf(s)
		if !ok {
			return nil, false
		}
	}
	return root, true
}

func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		return nil, false
	}
	res, ok := n.children[path]
	return res, ok
}

// newRouter return a new route
func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}
