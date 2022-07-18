package gyy

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// 封装对象
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求信息
	Path   string // 请求路由
	Method string
	Params map[string]string // 路由参数
	// 返回信息
	StatusCode int
	// 中间件
	handlers []HandlerFunc // 存储注册的中间件
	index    int           // 记录当前执行中间件顺序
	//
	engine *Engine // 存储 engine 实例
}

// 初始化 context
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}

// 执行下一个注册的中间件
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	// 如果中间件中不手动执行 Next()，由于 for 中的 c.index++，也能执行完后续的中间件
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// 中断执行
func (c *Context) Fail(code int, err string) {
	// 赋为最大值，取消后续中间件执行
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// 根据 key 获取路由参数
// 例如路由为："/abc/:name"，请求为 "/abc/John" c.Param("name") = John
func (c *Context) Param(key string) string {
	v, _ := c.Params[key]
	return v
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}
