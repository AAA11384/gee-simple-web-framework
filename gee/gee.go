package gee

import (
	"net/http"
)

type Engine struct {
	router *router
	RouterGroup
}

type HandlerFunc func(ctx *Context)

type HandlersChain []HandlerFunc

func New() *Engine {
	s := &Engine{
		router: newRouter(),
		RouterGroup: RouterGroup{
			root:     true,
			Handlers: nil,
			basePath: "/",
		},
	}
	s.RouterGroup.engine = s
	return s
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
	e.router.addRoute(method, pattern, f)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.router.handle(newContext(w, req))
}
