package middlewares

import (
	"github.com/gin-gonic/gin"
	"red_envelope/utils"
)

// RecoveryMiddleWare 捕获所有panic，并且返回错误信息
func RecoveryMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//c.JSON(http.StatusInternalServerError, gin.H{
				//	"code": utils.CODE_OTHER_ERROR,
				//	"msg":  utils.MSG_OTHER_ERROR,
				//})
				utils.MyLog.Error(err)
				return
			}
		}()
		c.Next()
	}
}
