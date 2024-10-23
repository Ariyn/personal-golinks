package main

import (
	"log"
	"net/http"
)

type Routing struct {
	path     string
	method   string
	function func(http.ResponseWriter, *http.Request) error
}

var routing = make([]*Routing, 0)

func addRouting(path string, method string, function func(http.ResponseWriter, *http.Request) error) {
	routing = append(routing, &Routing{
		path:     path,
		method:   method,
		function: function,
	})
}

func get(path string, function func(http.ResponseWriter, *http.Request) error) {
	addRouting(path, http.MethodGet, function)
}

func post(path string, function func(http.ResponseWriter, *http.Request) error) {
	addRouting(path, http.MethodPost, function)
}

func router(w http.ResponseWriter, r *http.Request) {
	pathFound := false

	for _, route := range routing {
		if r.URL.Path == route.path {
			pathFound = true

			if r.Method == route.method {
				err := route.function(w, r)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	}

	if pathFound {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	http.Error(w, "Not Found", http.StatusNotFound)
}

func main() {
	get("/hello", func(w http.ResponseWriter, r *http.Request) error {
		log.Println("Received request", r.URL.Path, r.Method, r.Body)
		w.Write([]byte("Hello, World!"))
		return nil
	})

	s := &http.Server{
		Addr:    ":80",
		Handler: http.HandlerFunc(router),
	}

	log.Fatal(s.ListenAndServe())
}
