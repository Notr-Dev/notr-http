package notrhttp

import (
	"fmt"
	"net/http"
	"time"
)

type Job struct {
	Name     string
	Interval time.Duration
	Job      func() error
}

type Server struct {
	Name    string
	Port    string
	Version string

	Middlewares []*Middleware
	Routes      []Route
	Services    []*Service
	Jobs        []Job
}

func NewServer(server Server) *Server {
	if server.Name == "" {
		server.Name = "Unnamed Server"
	}
	if server.Port == "" {
		panic("Port is required")
	}
	if server.Version == "" {
		panic("Version is required")
	}
	server.Routes = []Route{}
	server.Services = []*Service{}
	server.Jobs = []Job{}
	server.Middlewares = []*Middleware{}
	return &server
}

func (s *Server) RegisterService(service *Service) {
	s.Services = append(s.Services, service)
}
func (s *Server) RegisterJob(job Job) {
	s.Jobs = append(s.Jobs, job)
}
func (s *Server) RegisterMiddleware(middleware Middleware) {
	s.Middlewares = append(s.Middlewares, &middleware)
}

func (s *Server) Run() error {

	err := setupServices(s)
	if err != nil {
		return err
	}

	s.Get("/", func(rw Writer, r *Request) {
		rw.RespondWithSuccess(map[string]string{
			"message": fmt.Sprintf("Welcome to the %s Rest API.", s.Name),
			"version": s.Version,
		})
	})

	appendRoutes(s)

	startJobs(s)

	fmt.Println("Started " + s.Name)

	return http.ListenAndServe(s.Port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wasRightPath := false
		for _, route := range s.Routes {
			isMatch, params := matchPath(r.URL.Path, route.Path)
			if isMatch {
				// fmt.Printf("Matched %s with %s\n", r.URL.Path, route.Path)
				// fmt.Println(params)
				wasRightPath = true
				if r.Method == route.Method {
					fmt.Println("Matched route", r.URL.Path)
					route.Handler(Writer{w, false}, &Request{r, params})
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

func appendRoutes(s *Server) {
	for i := range s.Routes {
		route := &s.Routes[i]

		combinedMiddlewares := make([]*Middleware, 0, len(s.Middlewares)+len(route.Middlewares))
		combinedMiddlewares = append(combinedMiddlewares, s.Middlewares...)
		combinedMiddlewares = append(combinedMiddlewares, route.Middlewares...)

		route.Middlewares = combinedMiddlewares
	}

	for i := range s.Services {
		service := s.Services[i]
		for j := range service.Routes {
			route := &service.Routes[j]

			combinedMiddlewares := make([]*Middleware, 0, len(s.Middlewares)+len(route.Middlewares)+len(service.Middlewares))
			combinedMiddlewares = append(combinedMiddlewares, s.Middlewares...)
			combinedMiddlewares = append(combinedMiddlewares, service.Middlewares...)
			combinedMiddlewares = append(combinedMiddlewares, route.Middlewares...)

			route.Middlewares = combinedMiddlewares
		}
		s.Routes = append(s.Routes, service.Routes...)
	}

	for i := range s.Routes {
		originalHandler := s.Routes[i].Handler
		for j := len(s.Routes[i].Middlewares) - 1; j >= 0; j-- {
			originalHandler = (*s.Routes[i].Middlewares[j])(originalHandler)
		}
		s.Routes[i].Handler = originalHandler
	}
}

func setupServices(s *Server) error {
	if len(s.Services) == 0 {
		fmt.Println("No services to start")
		return nil
	}

	fmt.Printf("Starting %d services\n", len(s.Services))

	for {
		allInitialized := true
		for _, service := range s.Services {
			if !service.IsInitialized && service.CanRun() {
				allInitialized = false
				err := service.initialize(s)
				if err != nil {
					return err
				}
			}

			fmt.Printf("Service: %s, IsInit: %t, CanRun %t\n", service.Name, service.IsInitialized, service.CanRun())

			for _, dep := range service.Dependencies {
				fmt.Printf("Dependency: %s, IsInit: %t\n", dep.Name, dep.IsInitialized)
			}
		}
		if allInitialized {
			break
		}
		time.Sleep(5 * time.Second)
	}
	return nil
}

func startJobs(s *Server) {
	for _, job := range s.Jobs {
		go func(job Job) {
			for {
				err := job.Job()
				if err != nil {
					fmt.Printf("Error in job %s: %s\n", job.Name, err.Error())
				}
				time.Sleep(job.Interval)
			}
		}(job)
	}
}
