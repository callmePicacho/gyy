package main

import (
	"gyy/gyy"
	"log"
	"net/http"
	"time"
)

// 计算执行时间
func Logger() gyy.HandlerFunc {
	return func(c *gyy.Context) {
		t := time.Now()
		log.Println("Logger")
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func onlyForV2() gyy.HandlerFunc {
	return func(c *gyy.Context) {
		t := time.Now()
		log.Println("onlyForV2")
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gyy.New()

	r.Use(Logger())
	r.GET("/", func(c *gyy.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>\n")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *gyy.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":9999")
}
