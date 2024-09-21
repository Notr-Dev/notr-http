package notrhttp

import "fmt"

type Service[T any] struct {
	Name          string
	isInitialized bool
	initFunction  func(service *Service[T], server *Server[T]) error

	Dependencies []*Service[T]
}

func NewService[T any](name string) Service[T] {
	return Service[T]{
		Name:          name,
		isInitialized: false,
		initFunction:  func(service *Service[T], server *Server[T]) error { return nil },
	}
}

func (s *Service[T]) SetInitFunction(initFunction func(service *Service[T], server *Server[T]) error) {
	s.initFunction = initFunction
}

func (s *Service[T]) initialize(server *Server[T]) error {
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

func (s *Service[T]) AddDependency(dep *Service[T]) {
	if dep == s {
		panic(fmt.Sprintf("Service %s cannot have itself as a dependency", s.Name))
	}
	s.Dependencies = append(s.Dependencies, dep)
}

func (s *Service[T]) CanRun() bool {
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
