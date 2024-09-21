package notrhttp

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

	routes   []Route
	services []*Service
	jobs     []Job
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
	server.routes = []Route{}
	server.services = []*Service{}
	server.jobs = []Job{}
	return &server
}

func (s *Server) RegisterService(service *Service) {
	s.services = append(s.services, service)
}
func (s *Server) RegisterJob(job Job) {
	s.jobs = append(s.jobs, job)
}

func (s *Server) Run() error {

	if len(s.services) == 0 {
		fmt.Println("No services to start")
	} else {

		fmt.Printf("Starting %d services\n", len(s.services))

		for {
			allInitialized := true
			for _, service := range s.services {
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

	for _, service := range s.services {
		for _, route := range service.Routes {
			s.routes = append(s.routes, route)
		}
	}

	fmt.Println("Started " + s.Name)
	s.Get("/", func(rw Writer, r *Request) {
		rw.RespondWithSuccess(map[string]string{
			"message": fmt.Sprintf("Welcome to the %s Rest API.", s.Name),
			"version": s.Version,
		})
	})

	for _, job := range s.jobs {
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
		for _, route := range s.routes {
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
func matchPath(path string, pattern string) (is bool, params map[string]string) {
	fmt.Println("Path: ", path)
	fmt.Println("Pattern: ", pattern)
	if path == "" {
		path = "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	pathComponents := strings.Split(path, "/")
	patternComponents := strings.Split(pattern, "/")

	fmt.Println("Path Components: ", pathComponents)
	fmt.Println("Pattern Components: ", patternComponents)

	if len(pathComponents) != len(patternComponents) {
		return false, map[string]string{}
	}

	params = map[string]string{}

	for i, component := range patternComponents {
		if component != pathComponents[i] && !strings.HasPrefix(component, "{") {
			return false, map[string]string{}
		} else {
			if strings.HasPrefix(component, "{") {
				params[component[1:len(component)-1]] = pathComponents[i]
			}
		}
	}
	return true, params
}
func (s *Server) genericHandler(method string, path string, handler Handler) {
	s.routes = append(s.routes,
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

func (s *Server) StaticServe(path string, dir string) {
	s.routes = append(s.routes,
		Route{
			Method: http.MethodGet,
			Path:   path + "/{filename}",
			Handler: func(rw Writer, r *Request) {
				filename := r.Params["filename"]
				filePath := filepath.Join(dir, filename)
				fmt.Println("Serving file: ", filePath)

				file, err := os.Open(filePath)
				if err != nil {
					http.Error(rw, "File not found", http.StatusNotFound)
					return
				}
				defer file.Close()

				rw.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(filePath)))

				if _, err := io.Copy(rw, file); err != nil {
					http.Error(rw, "Error serving file", http.StatusInternalServerError)
				}
			},
		},
	)
}
