package notrhttp

import "fmt"

type Service struct {
	Name    string
	Subpath string
	isReady bool

	Handlers     []Handler
	Dependencies []*Service
}

func (s *Service) NewService(name string, subpath string) *Service {
	return &Service{
		Name:    name,
		Subpath: subpath,
	}
}

func (s *Service) AddDependency(dep *Service) {
	if dep != s {
		panic(fmt.Sprintf("Service %s cannot have itself as a dependency", s.Name))
	}
	s.Dependencies = append(s.Dependencies, dep)
}

func (s *Service) CanRun() bool {
	if len(s.Dependencies) == 0 {
		return s.isReady
	}
	for _, dep := range s.Dependencies {
		if dep == s {
			panic(fmt.Sprintf("Service %s cannot have itself as a dependency", s.Name))
		}
		if !dep.CanRun() {
			return false
		}
	}
	return true
}
