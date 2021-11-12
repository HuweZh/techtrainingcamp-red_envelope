package main

import (
	"github.com/gin-gonic/gin"
	"huhusw.com/red_envelope/routers"
)

func main() {
	//创建一个新的路由引擎
	r := gin.Default()

	//加载路由
	routers.RoutersInit(r)

	//启动路由
	r.Run()

}
