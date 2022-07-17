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

	r.GET("/hello/:name", func(c *gyy.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *gyy.Context) {
		c.JSON(http.StatusOK, gyy.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")
}
