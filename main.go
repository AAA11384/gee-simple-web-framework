package main

import (
	"fmt"
	"gee/gee"
	"log"
	"net/http"
	"time"
)

func main() {
	r := gee.New()

	r.Static("/test/staticSources", "./gee/static")

	//基本访问验证
	r.GET("/test01", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>test 1 success</h1>")
	})
	//参数解析验证
	r.GET("/test02/:param", func(c *gee.Context) {
		c.String(http.StatusOK, "test 2 success, param is %s", c.Param("param"))
	})
	//表单参数解析
	r.POST("/test03", func(c *gee.Context) {
		c.String(http.StatusOK, "test 3 success, postParam is %s", c.PostForm("postKey"))
	})

	//路由组分组功能演示
	v1 := r.Group("/v1")
	{
		//返回JSON功能演示
		v1.GET("/test04", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"test04Key": "test04Value",
			})
		})
	}

	//注册中间件功能演示
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

	//panic恢复演示
	v3 := r.Group("/v3")
	{
		v3.GET("/test06", func(c *gee.Context) {
			arr := []int{1, 2, 3}
			fmt.Println(arr[4])
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":9999")
}
