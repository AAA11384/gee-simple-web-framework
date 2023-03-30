package main

import (
	"gee/gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.Static("/text/asd/abc", "./gee/static")
	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/v1")

	v1.Use(Logger(), Logger2())

	{
		v1.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}
	r.Run(":9999")
}

func Logger() gee.HandlerFunc {
	return func(c *gee.Context) {
		c.String(http.StatusOK, "123")
		c.Next()
		c.String(http.StatusOK, "10")
	}
}
func Logger2() gee.HandlerFunc {
	return func(c *gee.Context) {
		c.String(http.StatusOK, "456")
		c.Abort()
		c.String(http.StatusOK, "789")
	}
}
