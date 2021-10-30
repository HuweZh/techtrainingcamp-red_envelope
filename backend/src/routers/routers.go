package routers

import (
	"github.com/gin-gonic/gin"
	"huhusw.com/red_envelope/controller"
	"huhusw.com/red_envelope/middleware"
)

func RoutersInit(r *gin.Engine) {
	//拆红包路由
	r.POST("/snatch", middleware.InitMiddle, controller.RedEnvelopeController{}.Snatch)

	//拆红包路由
	r.POST("/open", middleware.InitMiddle, controller.RedEnvelopeController{}.Open)

	//拆红包路由
	r.POST("/get_wallet_list", middleware.InitMiddle, controller.RedEnvelopeController{}.GetWalletList)

	//测试数据库
	// r.GET("/test", controller.RedEnvelopeController{}.Test)
}
