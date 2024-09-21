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

	Routes   []Route
	Services []*Service
	Jobs     []Job
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
	return &server
}

func (s *Server) RegisterService(service *Service) {
	s.Services = append(s.Services, service)
}
func (s *Server) RegisterJob(job Job) {
	s.Jobs = append(s.Jobs, job)
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

	return http.ListenAndServe(s.Port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wasRightPath := false
		for _, route := range s.Routes {
			isMatch, params := matchPath(r.URL.Path, route.Path)
			if isMatch {
				fmt.Println("Path matched: ", route.Path)
				wasRightPath = true
				if r.Method == route.Method {
					route.Handler(Writer{w, false}, &Request{r, params})
					return
				}
			} else {
				fmt.Println("Path did not match: ", route.Path)
			}

		}
		if wasRightPath {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.NotFound(w, r)
	}))
}