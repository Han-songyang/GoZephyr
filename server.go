package GoZephyr

import "net/http"

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
}

var _ Server = &HTTPServer{}

type HTTPServer struct {
}

func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	s.serve(ctx)
}

// Start 启动服务器
func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) Post(path string, handler HandleFunc) {
}

func (s *HTTPServer) Get(path string, handler HandleFunc) {
}

func (s *HTTPServer) serve(ctx *Context) {

}
