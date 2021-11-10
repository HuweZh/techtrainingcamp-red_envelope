package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"red_envelope/utils"
)

type AdminController struct{}

type addMoneyRequest struct {
	Money int64 `json:"money"` //增加的资金
	Size  int64 `json:"size"`  //增加的红包数
}

func (con AdminController) AddMoneyHandler(c *gin.Context) {
	var request addMoneyRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.CODE_BINDJSON_ERROR,
			"msg":  utils.MSG_BINDJSON_ERROR,
		})
		panic(err)
	}
	err = utils.AddMoney(request.Money, request.Size)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.CODE_ADD_MONEY_ERROR,
			"msg":  utils.MSG_ADD_MONEY_ERROR,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": utils.CODE_SUCCESS,
		"msg":  utils.MSG_SUCCESS,
	})
	return
}
