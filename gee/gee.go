package gee

import (
	"html/template"
	"net/http"
	"path"
)

type Engine struct {
	router *router
	RouterGroup
	routerGroups  []*RouterGroup
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
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
		routerGroups:  []*RouterGroup{},
		htmlTemplates: nil,
		funcMap:       map[string]any{},
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
	c := newContext(w, req)
	c.engine = e
	e.router.handle(c)
}

// createStaticHandler 将静态文件目录转化为Handler，当访问来临的时候，拿到filepath
// 在目录Handler中打开文件，使文件Handler转化为一个文件的Handler，向请求ResWriter中
// 写静态文件内容
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.basePath, relativePath)
	//使用http.StripPrefix将absolutePath为路径的Handler转发到root为路径的目录Handler
	//即类似于http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("/tmp"))))
	//生成的Handler是一个加了转发的目录Handler
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		//目录Handler -> 文件Handler
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		//此时是文件Handler，而不是目录Handler
		//逻辑:context找req.Url,匹配路径,符合则转发到root为路径的Handler
		fileServer.ServeHTTP(c.W, c.Req)
	}
}

// Static 让root中的静态资源映射到relativePath上，访问静态资源只需在relativePath后加
// filepath即可
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}
