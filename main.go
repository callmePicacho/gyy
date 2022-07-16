package main

import (
	"gyy/gyy"
	"net/http"
)

func main() {
	r := gyy.New()

	r.GET("/", func(c *gyy.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>\n")
	})

	r.GET("/hello", func(c *gyy.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gyy.Context) {
		c.JSON(http.StatusOK, gyy.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}
