package routers

import (
	"github.com/gin-gonic/gin"
	"red_envelope/controllers/admin"
)

// AdminRoutersInit 首字母大写表示公有权限
func AdminRoutersInit(r *gin.Engine) {
	//路由分组
	apiRouters := r.Group("/admin")
	{
		apiRouters.POST("/add_money", admin.AdminController{}.AddMoneyHandler) //抢红包(加括号表示执行方法,{}表示实例化结构体)
	}
}
