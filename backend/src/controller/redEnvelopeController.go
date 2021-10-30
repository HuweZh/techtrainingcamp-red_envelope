package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"huhusw.com/red_envelope/entity"
)

/*
红包相关接口控制器
*/
type RedEnvelopeController struct {
}

//抢红包业务逻辑
func (con RedEnvelopeController) Snatch(c *gin.Context) {
	//测试数据库
	TestDB()

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
}

//打开红包业务逻辑
func (con RedEnvelopeController) Open(c *gin.Context) {
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
}

//获取钱包列表业务逻辑
func (con RedEnvelopeController) GetWalletList(c *gin.Context) {
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
	//中间件通信
	fmt.Println(c.Get("name"))
	//数据返回
	c.JSON(0, map[string]interface{}{
		"uid": user.Uid,
	})
}

//测试数据库
func (con RedEnvelopeController) Test(c *gin.Context) {
	TestDB()
}
