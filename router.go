package golinks

import (
	"fmt"
	"net/http"
)

type Routing struct {
	Path     string
	Method   string
	Function func(http.ResponseWriter, *http.Request) error
}

type Router struct {
	Routings []*Routing
}

func NewRouter() *Router {
	return &Router{
		Routings: make([]*Routing, 0),
	}
}

func (r *Router) AddRouting(path string, method string, function func(http.ResponseWriter, *http.Request) error) (err error) {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if path[0] != '/' {
		path = "/" + path
	}

	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	r.Routings = append(r.Routings, &Routing{
		Path:     path,
		Method:   method,
		Function: function,
	})

	return nil
}

func (r *Router) Get(path string, function func(http.ResponseWriter, *http.Request) error) (err error) {
	return r.AddRouting(path, http.MethodGet, function)
}

func (r *Router) Post(path string, function func(http.ResponseWriter, *http.Request) error) (err error) {
	return r.AddRouting(path, http.MethodPost, function)
}

func (r *Router) Redirect(path string, url string) (err error) {
	return r.Get(path, func(writer http.ResponseWriter, request *http.Request) error {
		http.Redirect(writer, request, url, http.StatusSeeOther)
		return nil
	})
}

func (r *Router) Route(writer http.ResponseWriter, request *http.Request) {
	pathFound := false

	method := request.Method
	path := request.URL.Path
	for _, route := range r.Routings {
		if path == route.Path {
			pathFound = true

			if method == route.Method {
				err := route.Function(writer, request)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	}

	if pathFound {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	http.Error(writer, "Not Found", http.StatusNotFound)
}
