package notrhttp

import (
	"fmt"
	"net/http"
)

type Server struct {
	Router *Router
	Name   string
	Port   string
}

func NewServer(port string) *Server {
	if port[0] != ':' {
		port = ":" + port
	}
	return &Server{
		Router: NewRouter(),
		Port:   port,
		Name:   "Unnamed Server",
	}
}

func (s *Server) Run() error {
	fmt.Println("Started server " + s.Name)
	return http.ListenAndServe(s.Port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Router.Mux.ServeHTTP(w, r)
	}))
}

func (s *Server) SetName(name string) {
	s.Name = name
}

type Router struct {
	Mux http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		Mux: *http.NewServeMux(),
	}
}
