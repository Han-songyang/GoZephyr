package GoZephyr

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

// AddRoute is a method to Create a node, including paths and HandleFunc
// and register it in the routing tree.
func (r *router) AddRoute(method string, path string, handler HandleFunc) {

}

// newRouter return a new route
func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}
