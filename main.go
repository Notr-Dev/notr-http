package notrhttp

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	Routes   []Route
	Name     string
	Port     string
	Version  string
	Services []*Service
}

func NewServer(port string, version string) *Server {
	if port[0] != ':' {
		port = ":" + port
	}
	return &Server{
		Routes:   []Route{},
		Port:     port,
		Name:     "Unnamed Server",
		Version:  version,
		Services: make([]*Service, 0),
	}
}

func (s *Server) RegisterService(service *Service) {
	s.Services = append(s.Services, service)
}

func (s *Server) Run() error {

	fmt.Printf("Starting %d services\n", len(s.Services))

	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(s.Services))

	for _, service := range s.Services {
		wg.Add(1)
		go func(service *Service) {
			defer wg.Done()
			for !service.CanRun() {
				// Wait until the service can run
				// You might want to add a sleep here to avoid busy waiting
				time.Sleep(time.Second)
			}
			if !service.isInitialized {
				fmt.Printf("Initializing %s service\n", service.Name)
				err := service.initialize()
				if err != nil {
					errChan <- err
					return
				}
				mu.Lock()
				service.isInitialized = true
				fmt.Printf("Service %s initialized\n", service.Name)
				mu.Unlock()
			}
		}(service)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	fmt.Println("Services status:")
	for _, service := range s.Services {
		status := "not initialized"
		if service.isInitialized {
			status = "initialized"
		}
		fmt.Printf("Service: %s, Status: %s\n", service.Name, status)
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

func (s *Server) SetName(name string) {
	s.Name = name
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
