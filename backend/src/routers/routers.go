package routers

import (
	"net/http"

	"red_envelope/src/entity"
	"github.com/gin-gonic/gin"
)

func RoutersInit(r *gin.Engine) {
	//拆红包路由
	r.POST("/snatch", func(c *gin.Context) {
		//用户实体
		user := entity.User{}
		//接收post的json数据
		// err := c.ShouldBindBodyWith(&user, binding.JSON)
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()})
			return
		}
		//数据返回
		c.JSON(0, map[string]interface{}{
			"uid": user.Uid,
		})
	})

	//拆红包路由
	r.POST("/open", func(c *gin.Context) {
		//红包实体
		// envelope := entity.Envelope{}
		//用户实体
		// user := entity.User{}
		//接收post的json数据
		// err := c.ShouldBindBodyWith(&user, binding.JSON)
		m := map[string]interface{}{
			"uid":         1,
			"envelope_id": 1,
		}
		err := c.ShouldBindJSON(&m)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()})
			return
		}
		//数据返回
		c.JSON(0, map[string]interface{}{
			"uid":         m["uid"],
			"envelope_id": m["envelope_id"],
		})
	})

	//拆红包路由
	r.POST("/get_wallet_list", func(c *gin.Context) {
		//用户实体
		user := entity.User{}
		//接收post的json数据
		// err := c.ShouldBindBodyWith(&user, binding.JSON)
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()})
			return
		}
		//数据返回
		c.JSON(0, map[string]interface{}{
			"uid": user.Uid,
		})
	})
}
