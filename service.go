package notrhttp

import "fmt"

type Service struct {
	Name          string
	isInitialized bool
	initFunction  func(service *Service, server *Server) error

	Dependencies []*Service
}

func NewService(opts ...func(*Service)) *Service {
	service := &Service{
		Name:          "Unnamed Service",
		isInitialized: false,
		initFunction:  func(service *Service, server *Server) error { return nil },
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

func WithServiceInitFunction(initFunction func(service *Service, server *Server) error) func(*Service) {
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

func (s *Service) initialize(server *Server) error {
	if s.isInitialized {
		return fmt.Errorf("Service %s is already initialized", s.Name)
	}
	err := s.initFunction(s, server)
	if err != nil {
		return err
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
