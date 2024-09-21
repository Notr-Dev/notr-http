package notrhttp

import (
	"fmt"
	"net/http"
	"strings"
)

type Route struct {
	Method  string
	Path    string
	Handler Handler
}

type Router struct {
	Routes []Route
}

type Handler func(rw Writer, r *Request)

type Request struct {
	*http.Request
	Params map[string]string
}

type Writer struct {
	http.ResponseWriter
	HadRespondedEarlier bool
}

func (s *Server) genericHandler(method string, path string, handler Handler) {
	s.Routes = append(s.Routes, Route{Method: method, Path: path, Handler: handler})
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
