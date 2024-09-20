package notrhttp

import "net/http"

type Server struct {
	Router *Router
	port   string
}

func NewServer(port string) *Server {
	if port[0] != ':' {
		port = ":" + port
	}
	return &Server{
		Router: NewRouter(),
		port:   port,
	}
}

type Router struct {
	Mux http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		Mux: *http.NewServeMux(),
	}
}
