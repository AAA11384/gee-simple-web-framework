# GEE

## Intro

It is a practice program created a simple web service framework imitating gin. It costs me four days, the final day used to test and write `README.md`.

I referred a article [7天用Go从零实现Web框架Gee教程 | 极客兔兔 (geektutu.com)](https://geektutu.com/post/gee.html), but not all the same, I also learned source code from gin v1.18 at the same time. 

In my gee program, it realized some basic function: parameter matching, prefix tree, router, static sources, panic recover and the support of HTML template.

## Construction

There are four main part in this program: Engine, route, prefix tree and context.

**route**

the route struct has two fields, handler restore handler chain of a pattern in specific method, roots restore the root node of every methods.

```go
type router struct {
	handler map[string]HandlersChain
	roots   map[string]*node
}
```

The node of prefix tree has four fields, if there is `:` or `*` operator in path, is `isWild` field of counterpart node will be true. Only the pattern field of end node in a tree is not a blank string. When a request coming, we search counterpart node of this path in prefix tree, then use request method and pattern field to generate a key, so we can get `HandlersChain` in `router.handler`.At the same time, we also can generate the map of parameters and value.(For more details, see getRoute function in router.go)

```go
type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}
```

POST and GET function in `gee.go` will call `addRoute`  finally to add the handler and pattern in router. After compiling, all the information is prepared in router.

**Engine**

Generally speaking, one project has one engine, it contains all the information.

```go
type Engine struct {
	router *router
	RouterGroup
	routerGroups  []*RouterGroup
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}
```

**RouterGroup**

It contains `RouterGroup` struct, so it also can use the function of `RouterGroup`. a `RouterGroup` means a Group of router, we can set a middleware for one group or set pattern for a group, so that, the handlers will all be affected by it's group. We can get a RouterGroup from a Engine or another RouterGroup, the middlewares and pattern will be shared. 

Only root field in `Engine.RouterGroup` is true, others all false.

```go
type RouterGroup struct {
	root     bool
	Handlers HandlersChain
	basePath string
	engine   *Engine
}
```

**Context**

Every request has it own context, it contains basic information of a request. The program will depends on basic information search parameter map and HandlerChain of a request.

```go
type Context struct {
	//origin objects
	W   http.ResponseWriter
	Req *http.Request
	//request information
	Path   string
	Method string
	Params map[string]string
	//response info
	StatusCode int
	//middlewares
	handlers []HandlerFunc
	index    int
	//engine pointer
	engine *Engine
}
```

## Examples

#### 0.Static resources

```go
r := gee.New()

r.Static("/test/staticSources", "./gee/static")
```

`path `: `http://127.0.0.1:9999/test/staticSources/example.html`

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Hello Gee</title>
</head>
<body>
<h2>visit static source success !!!</h2>
</body>
</html>
```

#### 1.basic function

```go
	r.GET("/test01", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>test 1 success</h1>")
	})
```

`path` ：`http://127.0.0.1:9999/test01`

```html
<h1>test 1 success</h1>
```

#### 2.parameters praseing

```go
	r.GET("/test02/:param", func(c *gee.Context) {
		c.String(http.StatusOK, "test 2 success, param is %s", c.Param("param"))
	})
```

`path`:`http://127.0.0.1:9999/test02/aaa11384`

```
test 2 success, param is aaa11384
```

#### 3.post praseing

```go
	r.POST("/test03", func(c *gee.Context) {
		c.String(http.StatusOK, "test 3 success, postParam is %s", c.PostForm("postKey"))
	})
```

`use postman ` `path`:`127.0.0.1:9999/test03`

`postKey - postValue`

```
test 3 success, postParam is postValue
```

#### 4.route group

```go
	v1 := r.Group("/v1")
	{
		v1.GET("/test04", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"test04Key": "test04Value",
			})
		})
	}
```

`path`:`127.0.0.1:9999/test04`

```
{"test04Key":"test04Value"}
```

#### 5.middleware

```go
	v2 := r.Group("/v2")
	v2.Use(func(c *gee.Context) {
		log.Printf("a request coming at %s", time.Now().String())
	})
	{
		v2.POST("/test05", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}
```

`use postman ` `path`:`127.0.0.1:9999/v2/test05`

`username - aaa11384`

`password - 123456`

```json
{
    "password": "123456",
    "username": "aaa11384"
}
```

with`2023/03/30 19:08:22 a request coming at 2023-03-30 19:08:22.6810024 +0800 CST m=+52.969624301`  in console.

#### 6.panic recovery

```go
	v3 := r.Group("/v3")
	{
		v3.GET("/test06", func(c *gee.Context) {
			arr := []int{1, 2, 3}
			fmt.Println(arr[4])
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}
```

`path`:`127.0.0.1:9999/v3/test06`

with`2023/03/30 19:08:22 a request coming at 2023-03-30 19:08:22.6810024 +0800 CST m=+52.969624301
2023/03/30 19:11:12 http: panic serving 127.0.0.1:5780: runtime error: index out of range [4] with length 3` in console, and the program isn't exit.

