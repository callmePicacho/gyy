package main

import (
	"gyy/gyy"
	"net/http"
)

type student struct {
	Name string
	Age  int8
}

func main() {
	r := gyy.Default()

	r.GET("/", func(c *gyy.Context) {
		c.String(http.StatusOK, "hello lyy\n")
	})

	r.GET("/panic", func(c *gyy.Context) {
		arr := []string{"1"}
		c.String(http.StatusOK, arr[100]) // panic
	})

	r.Run(":9999")
}
