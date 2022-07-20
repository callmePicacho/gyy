# gyy
仿 Gin Go Web 框架

参考兔兔的实现：https://geektutu.com/post/gee.html

实现了：
1. 通过实现 ServerHttp 接口，接管 HTTP 请求
2. 实现 Context 封装 request 和 response，并提供对 JSON、HTML 等类型返回的接口
3. 通过实现 Trie 前缀树，实现通配符 * 和参数匹配 : 的动态路由功能
4. 通过前缀区分，实现对路由的分组控制
5. 设计实现用户可自定义中间件功能，并将路由处理函数也作为中间件的一部分
6. 通过动态路由实现静态服务器，借助 Go 内置系统库提供服务端渲染能力
7. 手动为 engine 添加错误恢复的中间件，防止由于 panic 系统宕机
