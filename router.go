package golinks

import "net/http"

type Routing struct {
	path     string
	method   string
	function func(http.ResponseWriter, *http.Request) error
}

type Router struct {
	Routings []*Routing
}

func (r *Router) AddRouting(path string, method string, function func(http.ResponseWriter, *http.Request) error) {
	r.Routings = append(r.Routings, &Routing{
		path:     path,
		method:   method,
		function: function,
	})
}

func (r *Router) Get(path string, function func(http.ResponseWriter, *http.Request) error) {
	r.AddRouting(path, http.MethodGet, function)
}

func (r *Router) post(path string, function func(http.ResponseWriter, *http.Request) error) {
	r.AddRouting(path, http.MethodPost, function)
}

func (r *Router) Route(writer http.ResponseWriter, request *http.Request) {
	pathFound := false

	method := request.Method
	path := request.URL.Path
	for _, route := range r.Routings {
		if path == route.path {
			pathFound = true

			if method == route.method {
				err := route.function(writer, request)
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
