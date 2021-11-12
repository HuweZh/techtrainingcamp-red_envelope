package routers

import (
	"github.com/gin-gonic/gin"
	"huhusw.com/red_envelope/controller"
	"huhusw.com/red_envelope/middleware"
)

func RoutersInit(r *gin.Engine) {
	//拆红包路由
	r.POST("/snatch", middleware.SnatchMiddle, controller.RedEnvelopeController{}.Snatch)

	//开红包路由
	r.POST("/open", middleware.OpenMiddle, controller.RedEnvelopeController{}.Open)

	//拆红包路由
	r.POST("/get_wallet_list", middleware.WalletMiddle, controller.RedEnvelopeController{}.GetWalletList)
}
