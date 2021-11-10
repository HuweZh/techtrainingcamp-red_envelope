package middlewares

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/gin-gonic/gin"
)

// ApiMiddleware 中间件，Group、GET、POST等方法中可以加入多个函数用于做中间件
type ApiMiddleware struct{}

// ApiMiddleware 中间件
func (con ApiMiddleware) ApiMiddleware(c *gin.Context) {
	e, err := sentinel.Entry("request")
	if err != nil {
		//c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		//	"code": utils.CODE_LIMIT_HAS_BEEN_REACHED_ERROR,
		//	"msg":  utils.MSG_LIMIT_HAS_BEEN_REACHED_ERROR,
		//})
		//utils.MyLog.Error(err)
		//return
	}
	e.Exit()
	c.Next()
}
