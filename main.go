package notrhttp

import (
	"fmt"
	"net/http"
)

type Server struct {
	Router  *http.ServeMux
	Name    string
	Port    string
	Version string
}

func NewServer(port string, version string) *Server {
	if port[0] != ':' {
		port = ":" + port
	}
	return &Server{
		Router:  http.NewServeMux(),
		Port:    port,
		Name:    "Unnamed Server",
		Version: version,
	}
}

func (s *Server) Run() error {
	fmt.Println("Started " + s.Name)
	s.Get("/", func(rw Writer, r *Request) {
		rw.RespondWithSuccess(map[string]string{
			"message": fmt.Sprintf("Welcome to the %s Rest API.", s.Name),
			"version": s.Version,
		})
	})
	s.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	return http.ListenAndServe(s.Port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Router.ServeHTTP(w, r)
	}))
}

func (s *Server) SetName(name string) {
	s.Name = name
}

type Request struct {
	*http.Request
}

type Handler func(rw Writer, r *Request)

func (s *Server) genericHandler(method string, path string, handler Handler) {
	s.Router.HandleFunc(method+" "+path, func(w http.ResponseWriter, r *http.Request) {
		handler(Writer{w}, &Request{r})
	})
}

func (s *Server) Get(path string, handler Handler) {
	s.genericHandler(http.MethodGet, path, handler)
}

func (s *Server) Post(path string, handler Handler) {
	s.genericHandler(http.MethodPost, path, handler)
}

func (s *Server) Put(path string, handler Handler) {
	s.genericHandler(http.MethodPut, path, handler)
}

func (s *Server) Delete(path string, handler Handler) {
	s.genericHandler(http.MethodDelete, path, handler)
}

func (s *Server) Patch(path string, handler Handler) {
	s.genericHandler(http.MethodPatch, path, handler)
}
