package controller

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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
	//1.获取请求参数
	para := getParamter(c)

	//2.根据id获取用户信息
	var user models.User
	data, err := commons.GetRedis().Get(c, strconv.Itoa(para.Uid)).Result()
	if err != nil {
		user := models.GetUser(para.Uid)
		//超时时间的单位为微秒，100*1000000000 是100秒
		commons.GetRedis().Set(c, strconv.Itoa(para.Uid), user, 100*1000000000)
	} else {
		//5.TODO 更新缓存
		json.Unmarshal([]byte(data), &user)
		commons.GetRedis().Expire(c, strconv.Itoa(para.Uid), 100*1000000000)
		fmt.Println(user)
	}

	//3.TODO 获取一个红包id
	// rand.Seed(time.Now().UnixNano())
	// envelopeId := rand.Intn(100000) + 1000
	envelopeId := commons.GetID()

	//4.TODO 存储此红包为该用户的一个红包，并且更新cur_count
	var envelope models.Envelope

	envelope.EnvelopeId = envelopeId
	envelope.UserId = para.Uid
	envelope.Opened = 0
	envelope.Value = 50
	envelope.SnatchTime = int(time.Now().Unix())
	models.SaveEnvelope(envelope)
	// models.UpdateCurCount(para.Uid, user.CurCount+1)

	//6.构建返回的数据
	d := gin.H{
		"envelope_id": envelopeId,
		"max_count":   user.MaxCount,
		"cur_count":   user.CurCount + 1,
	}
	r(c, d)
}

//打开红包业务逻辑
func (con RedEnvelopeController) Open(c *gin.Context) {
	//1.获取请求参数
	para := getParamter(c)
	//2.根据红包id获取红包信息
	envelope := models.GetEnvelope(para.Envelope_id)

	//4.TODO 更新红包的状态
	models.UpdateState(para.Envelope_id, 1)
	//5.TODO 更新缓存

	//6.构建返回的数据
	data := gin.H{
		"value": envelope.Value,
	}
	r(c, data)
}

//获取钱包列表业务逻辑
func (con RedEnvelopeController) GetWalletList(c *gin.Context) {
	//1.获取请求参数
	para := getParamter(c)
	//2.根据用户id获取当前的红包信息
	envelopes := models.GetEnvelopeList(para.Uid)
	//4.TODO 更新缓存

	//5.构建返回的数据
	var amount = 0
	for _, value := range envelopes {
		amount += value.Value
	}
	data := gin.H{
		"amount":        amount,
		"envelope_list": envelopes,
	}
	r(c, data)
}

func testDB(c *gin.Context) {
	//数据库测试
	user := []models.User{}
	commons.GetDB().Find(&user)
	for index, value := range user {
		fmt.Println(index, "  ", value)
	}
	//redis测试
	val, e := commons.GetRedis().Get(c, "b").Result()
	if e != nil {
		panic(e)
	}
	fmt.Println("b:", val)
	//中间件通信
	fmt.Println(c.Get("name"))
}

func getParamter(c *gin.Context) RequestParamter {
	//接受请求参数
	para := RequestParamter{}
	// err := c.ShouldBindBodyWith(&user, binding.JSON)
	err := c.ShouldBindJSON(&para)
	//请求参数错误
	if err != nil {
		fmt.Printf("Open file failed [Err:%s]\n", err.Error())
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
