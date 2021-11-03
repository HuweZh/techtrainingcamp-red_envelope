package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"huhusw.com/red_envelope/commons"
	"huhusw.com/red_envelope/models"
)

/*
红包相关接口控制器
*/
type RedEnvelopeController struct {
}

type RequestParamter struct {
	Uid         int `json:"uid"`
	Envelope_id int `json:"envelope_id"`
}

//抢红包业务逻辑
func (con RedEnvelopeController) Snatch(c *gin.Context) {
	//获取请求中的数据
	value, _ := c.Get("user")
	user := value.(models.User)

	// 存储此红包为该用户的一个红包
	var envelope models.Envelope = models.GetEnve(user.UserId)

	//超时时间的单位为微秒，100*1000000000 是100秒
	// commons.GetRedis().LSet(c, "uid"+strconv.Itoa(para.Uid), envelope)
	commons.GetRedis().RPush(c, "uid"+strconv.Itoa(user.UserId), envelope)
	commons.GetRedis().Expire(c, "uid"+strconv.Itoa(user.UserId), 600*1000000000)

	//构建返回的数据
	data := gin.H{
		"envelope_id": envelope.EnvelopeId,
		"max_count":   user.MaxCount,
		"cur_count":   user.CurCount,
	}

	//返回数据
	commons.R(c, commons.OK, commons.SUCCESS, data)
}

//打开红包业务逻辑
func (con RedEnvelopeController) Open(c *gin.Context) {
	//获取请求参数
	value, _ := c.Get("envelope")
	envelope := value.(models.Envelope)

	//构建返回的数据
	data := gin.H{
		"value": envelope.Value,
	}
	commons.R(c, commons.OK, commons.SUCCESS, data)
}

//获取钱包列表业务逻辑
func (con RedEnvelopeController) GetWalletList(c *gin.Context) {
	//获取请求携带的参数
	value, _ := c.Get("envelopes")

	//计算钱包的总数
	var amount = 0
	for _, value := range value.([]models.Envelope) {
		amount += value.Value
	}

	//构建返回数据，并返回
	data := gin.H{
		"amount":        amount,
		"envelope_list": value,
	}
	commons.R(c, commons.OK, commons.SUCCESS, data)
}
