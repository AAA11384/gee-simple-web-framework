package gee

type RouterGroup struct {
	root     bool
	Handlers HandlersChain
	basePath string
	engine   *Engine
}

func (group *RouterGroup) Group(relativePath string, f ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers: group.combineHandlers(f),
		basePath: group.calculateAbsolutePath(relativePath),
		engine:   group.engine,
	}
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	if group.basePath == "/" {
		return relativePath
	} else {
		return group.basePath + relativePath
	}
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.Handlers = append(group.Handlers, middlewares...)
}

func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	mergedHandler := make(HandlersChain, finalSize)
	return mergedHandler
}

func (group *RouterGroup) Handle(method, pattern string, f ...HandlerFunc) {
	HandlerChain := append(group.Handlers, f...)
	realPath := group.calculateAbsolutePath(pattern)
	group.engine.router.addRoute(method, realPath, HandlerChain...)
}

func (group *RouterGroup) POST(pattern string, f ...HandlerFunc) {
	group.Handle("POST", pattern, f...)
}

func (group *RouterGroup) GET(pattern string, f ...HandlerFunc) {
	group.Handle("GET", pattern, f...)
}
