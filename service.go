package notrhttp

import (
	"fmt"
)

type Service struct {
	Name          string
	Path          string
	isInitialized bool
	initFunction  func(service *Service) error

	Routes       []Route
	Dependencies []*Service
}

func NewService(opts ...func(*Service)) *Service {
	service := &Service{
		Name:          "Unnamed Service",
		isInitialized: false,
		initFunction:  func(service *Service) error { return nil },
		Dependencies:  []*Service{},
	}
	for _, opt := range opts {
		opt(service)
	}
	return service
}

func WithServiceName(name string) func(*Service) {
	return func(s *Service) {
		s.Name = name
	}
}

func WithServiceInitFunction(initFunction func(service *Service) error) func(*Service) {
	return func(s *Service) {
		s.initFunction = initFunction
	}
}

func WithServiceDependencies(dependencies ...*Service) func(*Service) {
	return func(s *Service) {
		for _, dep := range dependencies {
			if dep == s {
				panic("Service cannot depend on itself")
			}
			if dep == nil {
				panic("Dependency cannot be nil")
			}
		}
		s.Dependencies = dependencies
	}
}

func WithServiceSubpath(path string) func(*Service) {
	return func(s *Service) {
		s.Path = path
	}
}

func WithServiceRoutes(routes ...Route) func(*Service) {
	return func(s *Service) {
		s.Routes = routes
	}
}

func (s *Service) initialize() error {
	if s.isInitialized {
		return fmt.Errorf("Service %s is already initialized", s.Name)
	}
	err := s.initFunction(s)
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

	s.isInitialized = true
	return nil
}

func (s *Service) CanRun() bool {
	if len(s.Dependencies) == 0 {
		return true
	}
	for _, dep := range s.Dependencies {
		if !dep.isInitialized {
			return false
		}
	}
	return true
}
