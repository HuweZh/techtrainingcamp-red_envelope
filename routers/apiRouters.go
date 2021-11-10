package routers

import (
	"github.com/gin-gonic/gin"
	"red_envelope/controllers/api"
	"red_envelope/middlewares"
)

// ApiRoutersInit 首字母大写表示公有权限
func ApiRoutersInit(r *gin.Engine) {
	//路由分组
	apiRouters := r.Group("/api", middlewares.ApiMiddleware{}.ApiMiddleware)
	{
		apiRouters.POST("/snatch", api.ApiController{}.SnatchHandler) //抢红包(加括号表示执行方法,{}表示实例化结构体)
		apiRouters.POST("/open", api.ApiController{}.OpenHandler) //拆红包
		apiRouters.POST("/get_wallet_list", api.ApiController{}.GetWalletListHandler) //钱包列表
	}
}
