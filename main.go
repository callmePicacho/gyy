package main

import (
	"fmt"
	"gyy/gyy"
	"net/http"
)

func main() {
	r := gyy.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "我是/\n")
	})

	r.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "我是 hello\n")
	})

	r.Run(":9999")
}
