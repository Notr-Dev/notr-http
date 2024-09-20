package notrhttp

import (
	"fmt"
	"net/http"
)

type Server struct {
	Router *http.ServeMux
	Name   string
	Port   string
}

func NewServer(port string) *Server {
	if port[0] != ':' {
		port = ":" + port
	}
	return &Server{
		Router: http.NewServeMux(),
		Port:   port,
		Name:   "Unnamed Server",
	}
}

func (s *Server) Run() error {
	fmt.Println("Started " + s.Name)
	return http.ListenAndServe(s.Port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Router.ServeHTTP(w, r)
	}))
}

func (s *Server) SetName(name string) {
	s.Name = name
}

type CustomHandler func(rw ResponseWriterWrapper, r *http.Request)

func (s *Server) genericHandler(method string, path string, handler CustomHandler) {
	s.Router.HandleFunc(method+" "+path, func(w http.ResponseWriter, r *http.Request) {
		handler(ResponseWriterWrapper{w}, r)
	})
}

func (s *Server) Get(path string, handler CustomHandler) {
	s.genericHandler(http.MethodGet, path, handler)
}

func (s *Server) Post(path string, handler CustomHandler) {
	s.genericHandler(http.MethodPost, path, handler)
}

func (s *Server) Put(path string, handler CustomHandler) {
	s.genericHandler(http.MethodPut, path, handler)
}

func (s *Server) Delete(path string, handler CustomHandler) {
	s.genericHandler(http.MethodDelete, path, handler)
}

func (s *Server) Patch(path string, handler CustomHandler) {
	s.genericHandler(http.MethodPatch, path, handler)
}
