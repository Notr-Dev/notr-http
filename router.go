package notrhttp

import (
	"net/http"
	"strings"
)

type Route struct {
	Method      string        `json:"method"`
	Path        string        `json:"path"`
	Handler     Handler       `json:"-"`
	Middlewares []*Middleware `json:"-"`
}

type Router struct {
	Routes []Route
}

type Handler func(rw Writer, r *Request)

type Middleware func(next Handler) Handler

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
	if path == "" {
		path = "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	pathComponents := strings.Split(path, "/")
	patternComponents := strings.Split(pattern, "/")

	if len(pathComponents) != len(patternComponents) && !strings.Contains(pattern, "...") {
		return false, map[string]string{}
	}

	params = map[string]string{}

	for i, component := range patternComponents {
		if component != pathComponents[i] && !strings.HasPrefix(component, "{") {
			return false, map[string]string{}
		} else {
			if strings.HasPrefix(component, "{") && strings.HasSuffix(component, "}") {
				if strings.HasSuffix(component, "...}") {
					params[component[1:len(component)-4]] = strings.Join(pathComponents[i:], "/")
					break
				}
				params[component[1:len(component)-1]] = pathComponents[i]
			}
		}
	}

	return true, params
}
