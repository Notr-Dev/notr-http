package notrhttp

import "net/http"

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
}

type Writer struct {
	http.ResponseWriter
	HadRespondedEarlier bool
}
