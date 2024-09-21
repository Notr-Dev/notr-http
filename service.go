package notrhttp

import "fmt"

type Service struct {
	Name          string
	Subpath       string
	isReady       bool
	isInitialized bool
	initFunction  func(s *Service) error

	Handlers     []Handler
	Dependencies []*Service
}

func NewService(name string, subpath string) *Service {
	return &Service{
		Name:          name,
		Subpath:       subpath,
		isInitialized: false,
		initFunction:  func(s *Service) error { return nil },
	}
}

func (s *Service) SetInitFunction(initFunction func(s *Service) error) {
	s.initFunction = initFunction
}

func (s *Service) initialize() error {
	if s.isInitialized {
		return fmt.Errorf("Service %s is already initialized", s.Name)
	}
	err := s.initFunction(s)
	if err != nil {
		return err
	}
	s.isInitialized = true
	return nil
}

func (s *Service) AddDependency(dep *Service) {
	if dep == s {
		panic(fmt.Sprintf("Service %s cannot have itself as a dependency", s.Name))
	}
	s.Dependencies = append(s.Dependencies, dep)
}

func (s *Service) CanRun() bool {
	if len(s.Dependencies) == 0 {
		s.isReady = true
		return s.isReady
	}
	for _, dep := range s.Dependencies {
		if dep == s {
			panic(fmt.Sprintf("Service %s cannot have itself as a dependency", s.Name))
		}
		if !dep.isInitialized {
			s.isReady = false
			return false
		}
	}
	s.isReady = true
	return true
}
