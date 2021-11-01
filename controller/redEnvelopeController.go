package controller

import (
	"fmt"
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
	value, _ := c.Get("user")
	user := value.(models.User)
	//3.获取一个红包id
	// rand.Seed(time.Now().UnixNano())
	// envelopeId := rand.Intn(100000) + 1000
	// envelopeId := commons.GetID()

	//4.TODO 存储此红包为该用户的一个红包，并且更新cur_count
	var envelope models.Envelope = models.GetEnve(user.UserId)

	//超时时间的单位为微秒，100*1000000000 是100秒
	// commons.GetRedis().LSet(c, "uid"+strconv.Itoa(para.Uid), envelope)
	commons.GetRedis().RPush(c, "uid"+strconv.Itoa(user.UserId), envelope)
	commons.GetRedis().Expire(c, "uid"+strconv.Itoa(user.UserId), 600*1000000000)
	user.CurCount += 1
	var u models.UpdateData
	u.Type = models.INSERTENVELOPE
	u.Data = envelope
	//将数据传入写数据库的channel
	models.SetData(u)
	u.Type = models.UPDATEUSER
	u.Data = user
	models.SetData(u)
	commons.GetRedis().Set(c, strconv.Itoa(user.UserId), user, 100*1000000000)
	//6.构建返回的数据
	data := gin.H{
		"envelope_id": envelope.EnvelopeId,
		"max_count":   user.MaxCount,
		"cur_count":   user.CurCount,
	}
	//返回数据
	commons.R(c, 0, "success", data)
}

//打开红包业务逻辑
func (con RedEnvelopeController) Open(c *gin.Context) {
	//1.获取请求参数
	value, _ := c.Get("envelope")
	envelope := value.(models.Envelope)
	// fmt.Println("4 ", commons.GetRedis().LLen(c, "uid"+strconv.Itoa(para.Uid)))
	// fmt.Println(envelope)
	//先修改状态，再传入channel，在加入redis缓存
	envelope.Opened = 1
	var u models.UpdateData
	u.Type = models.UPDATEENVELOPESTATE
	u.Data = envelope
	//将数据传入写数据库的channel
	models.SetData(u)
	//超时时间的单位为微秒，100*1000000000 是100秒
	commons.GetRedis().RPush(c, "uid"+strconv.Itoa(envelope.UserId), envelope)
	commons.GetRedis().Expire(c, "uid"+strconv.Itoa(envelope.UserId), 600*1000000000)
	// envelope := models.GetEnvelope(para.Envelope_id)

	//4.TODO 更新红包的状态
	// models.UpdateState(para.Envelope_id, 1)
	//5.TODO 更新缓存

	//6.构建返回的数据
	data := gin.H{
		"value": envelope.Value,
	}
	commons.R(c, 0, "success", data)
}

//获取钱包列表业务逻辑
func (con RedEnvelopeController) GetWalletList(c *gin.Context) {

	value, _ := c.Get("envelopes")
	fmt.Println("value", value)
	// var envelopes []models.Envelope
	var uid int
	var amount = 0
	for _, value := range value.([]models.Envelope) {
		uid = value.UserId
		amount += value.Value
		fmt.Println("amount = ", amount)

	}
	//4.TODO 更新缓存
	commons.GetRedis().Expire(c, "uid"+strconv.Itoa(uid), 600*1000000000)
	// //5.构建返回的数据
	data := gin.H{
		"amount":        amount,
		"envelope_list": value,
	}
	commons.R(c, 0, "success", data)
}

func getParamter(c *gin.Context) RequestParamter {
	//接受请求参数
	para := RequestParamter{}
	// err := c.ShouldBindBodyWith(&user, binding.JSON)
	err := c.ShouldBindJSON(&para)
	//请求参数错误
	if err != nil {
		fmt.Printf("request paramter error [Err:%s]\n", err.Error())
		// c.AbortWithStatusJSON(
		// 	http.StatusInternalServerError,
		// 	gin.H{"error": err.Error()})
	}
	return para
}

func r(c *gin.Context, data map[string]interface{}) {
	c.JSON(0, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}
