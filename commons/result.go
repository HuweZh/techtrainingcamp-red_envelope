package commons

import "github.com/gin-gonic/gin"

func R(c *gin.Context, data map[string]interface{}) {
	c.JSON(0, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}
