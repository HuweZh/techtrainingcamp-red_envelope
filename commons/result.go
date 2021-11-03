package commons

import "github.com/gin-gonic/gin"

//构建返回的数据
func R(c *gin.Context, code int, msg string, data map[string]interface{}) {
	c.JSON(0, map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
