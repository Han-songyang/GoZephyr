package GoZephyr

import "net/http"

type HandleFunc func(ctx *Context)

type ICoreServer interface {
	http.Handler

	// Start an HTTP server
	Start(addr string) error

	// AddRoute is a method to Create a node, including paths and HandleFunc
	// and register it in the routing tree.
	addRoute(method string, path string, handler HandleFunc)
}

// Make sure struct implements this method
var _ ICoreServer = &CoreServer{}

type CoreServer struct {
	*router
}

// NewCoreServer Returns a HTTPServer
func NewCoreServer() *CoreServer {
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

func (s *CoreServer) Get(path string, handler HandleFunc) {
	s.addRoute(http.MethodGet, path, handler)
}

func (s *CoreServer) Post(path string, handler HandleFunc) {
	s.addRoute(http.MethodPost, path, handler)
}

func (s *CoreServer) serve(ctx *Context) {
	mi, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || mi.handler == nil {
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("Not Found"))
		return
	}
	mi.handler(ctx)
}
