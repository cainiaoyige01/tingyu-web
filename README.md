# tingyu-web

项⽬介绍：阅读net/http源码以及参考Gin框架的⼀些设置，做出来的
Web⼩型框架，对http标准库的进⾏封装。主要功能有注册路由、路由查
找、数据封装返回、分组控制、设置中间件以及错误恢复等。

项⽬难点：
实现⾃⼰路由器，定义路由规则弥补http标准库中没有动态路由功能，
如参数匹配：以及通配*这两个功能
实现分组路由控制以前缀进⾏区分并⽀持分组嵌套；实现错误处理机
制。
设计并实现Web框架的中间件机制，给框架安装⼀个插⼝，允许⽤⼾
实现⾃定义功能；
