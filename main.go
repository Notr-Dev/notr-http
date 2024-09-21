package notrhttp

import (
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	Name    string
	Port    string
	Version string

	Routes   []Route
	Services []*Service
}

func NewServer(options ...func(*Server)) *Server {
	svr := &Server{
		Name:    "Unnamed Server",
		Port:    ":8080",
		Version: "1.0.0",

		Routes:   []Route{},
		Services: []*Service{},
	}
	for _, o := range options {
		o(svr)
	}
	return svr
}

func WithServerPort(port string) func(*Server) {
	if port[0] != ':' {
		port = ":" + port
	}
	return func(s *Server) {
		s.Port = port
	}
}

func WithServerName(name string) func(*Server) {
	return func(s *Server) {
		s.Name = name
	}
}

func WithServerVersion(version string) func(*Server) {
	return func(s *Server) {
		s.Version = version
	}
}

func (s *Server) RegisterService(service *Service) {
	s.Services = append(s.Services, service)
}

func (s *Server) Run() error {

	if len(s.Services) == 0 {
		fmt.Println("No services to start")
	} else {

		fmt.Printf("Starting %d services\n", len(s.Services))

		for {
			allInitialized := true
			for _, service := range s.Services {
				if !service.isInitialized && service.CanRun() {
					allInitialized = false
					err := service.initialize()
					if err != nil {
						panic(err)
					}
				}

				fmt.Printf("Service: %s, IsInit: %t, CanRun %t\n", service.Name, service.isInitialized, service.CanRun())

				for _, dep := range service.Dependencies {
					fmt.Printf("Dependency: %s, IsInit: %t\n", dep.Name, dep.isInitialized)
				}
			}
			if allInitialized {
				break
			}
			time.Sleep(5 * time.Second)
		}

	}

	for _, service := range s.Services {
		for _, route := range service.Routes {
			s.Routes = append(s.Routes, route)
		}
	}

	fmt.Println("Started " + s.Name)
	s.Get("/", func(rw Writer, r *Request) {
		rw.RespondWithSuccess(map[string]string{
			"message": fmt.Sprintf("Welcome to the %s Rest API.", s.Name),
			"version": s.Version,
		})
	})
	return http.ListenAndServe(s.Port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wasRightPath := false
		for _, route := range s.Routes {
			if r.URL.Path == route.Path {
				wasRightPath = true
				if r.Method == route.Method {
					route.Handler(Writer{w, false}, &Request{r})
					return
				}
			}

		}
		if wasRightPath {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.NotFound(w, r)
	}))
}

func (s *Server) genericHandler(method string, path string, handler Handler) {
	s.Routes = append(s.Routes,
		Route{
			Method:  method,
			Path:    path,
			Handler: handler,
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
