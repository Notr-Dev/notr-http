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

func (s *Server) Run() error {
	return http.ListenAndServe(s.port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Router.Mux.ServeHTTP(w, r)
	}))
}

type Router struct {
	Mux http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		Mux: *http.NewServeMux(),
	}
}
