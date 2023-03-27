package gee

import "net/http"

type HandlerFunc func(ctx *Context)

type router struct {
	handler map[string]HandlerFunc
}

func newRouter() *router {
	return &router{map[string]HandlerFunc{}}
}

func (r *router) addRoute(method, pattern string, f HandlerFunc) {
	s := method + "-" + pattern
	r.handler[s] = f
}

func (r *router) handle(c *Context) {
	s := c.Method + "-" + c.Path
	if handler, ok := r.handler[s]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 Not Found :%s", c.Path)
	}
}
