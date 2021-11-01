package main

import (
	"github.com/gin-gonic/gin"
	"huhusw.com/red_envelope/routers"
)

func main() {
	//创建一个默认的路由引擎，里面默认加载了日志、错误抛出中间件
	r := gin.Default()

	//加载路由
	routers.RoutersInit(r)

	//启动路由
	r.Run()

}
