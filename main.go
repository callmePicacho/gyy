package main

import (
	"fmt"
	"gyy/gyy"
	"html/template"
	"log"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int8
}

// 计算执行时间
func Logger() gyy.HandlerFunc {
	return func(c *gyy.Context) {
		t := time.Now()
		log.Println("Logger")
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

// 格式化时间
func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gyy.New()
	r.Use(Logger())

	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "lyy", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/", func(c *gyy.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})

	r.GET("/students", func(c *gyy.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gyy.H{
			"title":  "gyy",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *gyy.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gyy.H{
			"title": "gyy",
			"now":   time.Date(2022, 7, 16, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999")
}
