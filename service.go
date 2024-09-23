package notrhttp

import (
	"fmt"
	"regexp"
)

type Service struct {
	PackageID     string                                       `json:"package_id"`
	Name          string                                       `json:"name"`
	Path          string                                       `json:"path"`
	IsInitialized bool                                         `json:"is_initialized"`
	InitFunction  func(service *Service, server *Server) error `json:"-"`

	Routes       []Route       `json:"routes"`
	Middlewares  []*Middleware `json:"middlewares"`
	Dependencies []*Service    `json:"dependencies"`
}

func NewService(service Service) *Service {
	if service.Name == "" {
		service.Name = "Unnamed Service"
	}
	if service.PackageID == "" {
		panic("Service must have a valid ID.")
	}
	if regexp.MustCompile(`\s`).MatchString(service.PackageID) {
		panic("Service ID cannot contain spaces.")
	}
	if len(service.Routes) > 0 {
		if service.Path == "" {
			panic("Path is required")
		}
		if service.Path[0] != '/' {
			panic("Path must start with a '/'")
		}
	}

	service.IsInitialized = false
	if service.InitFunction == nil {
		service.InitFunction = func(service *Service, server *Server) error { return nil }
	}
	if service.Dependencies == nil {
		service.Dependencies = []*Service{}
	}

	if len(service.Dependencies) > 0 {
		for _, dep := range service.Dependencies {
			if dep == nil {
				panic("Dependency cannot be nil")
			}
			if dep == &service {
				panic("Service cannot depend on itself")
			}
		}
	}

	if service.Routes == nil {
		service.Routes = []Route{}
	}
	return &service
}

func (s *Service) initialize(server *Server) error {
	if s.IsInitialized {
		return fmt.Errorf("Service %s is already initialized", s.Name)
	}
	err := s.InitFunction(s, server)
	if err != nil {
		return err
	}

	if len(s.Routes) > 0 {

		if s.Path[0] != '/' {
			panic("Service Path must start with a '/'")
		}

		for i, route := range s.Routes {
			if route.Path[0] != '/' {
				panic("Route Path must start with a '/'")
			}
			s.Routes[i].Path = s.Path + route.Path
		}
	}

	s.IsInitialized = true
	return nil
}

func (s *Service) CanRun() bool {
	if len(s.Dependencies) == 0 {
		return true
	}
	for _, dep := range s.Dependencies {
		if !dep.IsInitialized {
			return false
		}
	}
	return true
}
