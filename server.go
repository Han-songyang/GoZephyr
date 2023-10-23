package GoZephyr

import "net/http"

type HandleFunc func(ctx *Context)

type ICoreServer interface {
	http.Handler

	// Start an HTTP server
	Start(addr string) error

	// AddRoute is a method to Create a node, including paths and HandleFunc
	// and register it in the routing tree.
	AddRoute(method string, path string, handler HandleFunc)
}

// Make sure struct implements this method
var _ ICoreServer = &CoreServer{}

type CoreServer struct {
	*router
}

// NewHTTPServer Returns a HTTPServer
func NewHTTPServer() *CoreServer {
	return &CoreServer{
		router: newRouter(),
	}
}

// ServeHTTP Implementing methods in the Handler interface
func (s *CoreServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	s.serve(ctx)
}

// Start an HTTP server
func (s *CoreServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *CoreServer) Post(path string, handler HandleFunc) {
}

func (s *CoreServer) Get(path string, handler HandleFunc) {
}

func (s *CoreServer) serve(ctx *Context) {

}
