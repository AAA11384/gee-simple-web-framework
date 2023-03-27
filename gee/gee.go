package gee

import (
	"net/http"
)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
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
	e.router.handler[s] = f
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.router.handle(newContext(w, req))
}
