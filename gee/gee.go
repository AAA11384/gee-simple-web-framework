package gee

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	m map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{map[string]HandlerFunc{}}
}

func (e *Engine) GET(pattern string, f HandlerFunc) {
	e.addRoute("GET", pattern, f)
}

func (e *Engine) POST(pattern string, f HandlerFunc) {
	e.addRoute("POST", pattern, f)
}

func (e *Engine) Run(path string) {
	http.ListenAndServe(path, e)
}

func (e *Engine) addRoute(method, pattern string, f HandlerFunc) {
	s := method + "-" + pattern
	e.m[s] = f
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s := req.Method + "-" + req.URL.Path
	if handler, ok := e.m[s]; ok {
		handler(w, req)
	} else {
		fmt.Fprint(w, "404 Not Found")
	}
}
