package gee

import (
	"net/http"
	"strings"
)

type router struct {
	handler map[string]HandlersChain
	roots   map[string]*node
}

func newRouter() *router {
	return &router{
		handler: map[string]HandlersChain{},
		roots:   map[string]*node{},
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	res := make([]string, 0)
	for _, v := range vs {
		if v != "" {
			res = append(res, v)
			if v[0] == '*' {
				break
			}
		}
	}
	return res
}

func (r *router) addRoute(method, pattern string, f ...HandlerFunc) {
	s := method + "-" + pattern
	r.handler[s] = f

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = new(node)
	}
	r.roots[method].insert(pattern, parsePattern(pattern), 0)
}

func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for i, v := range parts {
			if v[0] == ':' {
				params[v[1:]] = searchParts[i]
			}
			//pattern为*123/456/789 传入a/b/c为path,则key为123,value为a/b/c
			if v[0] == '*' && len(v) > 1 {
				params[v[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, m := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = m
		s := c.Method + "-" + n.pattern
		c.handlers = r.handler[s]
		c.Next()
	} else {
		c.String(http.StatusNotFound, "404 Not Found")
	}
}
