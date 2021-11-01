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
		user = models.GetUser(para.Uid)
		//此用户的红包抢完了
		if user.MaxCount == user.CurCount {
			//返回数据
			r(c, nil)
		}
	} else {
		//5.更新用户的状态，并更新缓存
		json.Unmarshal([]byte(data), &user)
		//此用户的红包抢完了
		if user.MaxCount == user.CurCount {
			//返回数据
			r(c, nil)
		}
		user.CurCount += 1
	}
	//超时时间的单位为微秒，100*1000000000 是100秒
	commons.GetRedis().Set(c, strconv.Itoa(para.Uid), user, 100*1000000000)

	//3.获取一个红包id
	// rand.Seed(time.Now().UnixNano())
	// envelopeId := rand.Intn(100000) + 1000
	// envelopeId := commons.GetID()

	//4.TODO 存储此红包为该用户的一个红包，并且更新cur_count
	var envelope models.Envelope

	envelope.EnvelopeId = commons.GetID()
	envelope.UserId = para.Uid
	envelope.Opened = 0
	envelope.Value = 50
	envelope.SnatchTime = int(time.Now().Unix())

	//超时时间的单位为微秒，100*1000000000 是100秒
	// commons.GetRedis().LSet(c, "uid"+strconv.Itoa(para.Uid), envelope)
	commons.GetRedis().RPush(c, "uid"+strconv.Itoa(para.Uid), envelope)
	commons.GetRedis().Expire(c, "uid"+strconv.Itoa(para.Uid), 600*1000000000)

	var u models.UpdateData
	u.Type = models.INSERTENVELOPE
	u.Data = envelope
	//将数据传入写数据库的channel
	models.SetData(u)
	u.Type = models.UPDATEUSER
	u.Data = user
	models.SetData(u)

	//6.构建返回的数据
	d := gin.H{
		"envelope_id": envelope.EnvelopeId,
		"max_count":   user.MaxCount,
		"cur_count":   user.CurCount + 1,
	}
	//返回数据
	r(c, d)
}

//打开红包业务逻辑
func (con RedEnvelopeController) Open(c *gin.Context) {
	//1.获取请求参数
	para := getParamter(c)
	//2.根据红包id获取红包信息
	//尝试从redis中获取
	var envelope models.Envelope
	// fmt.Println("1 ", commons.GetRedis().LLen(c, "uid"+strconv.Itoa(para.Uid)))
	d, err := commons.GetRedis().LRange(c, "uid"+strconv.Itoa(para.Uid), 0, -1).Result()
	// fmt.Println("2 ", commons.GetRedis().LLen(c, "uid"+strconv.Itoa(para.Uid)))
	if err != nil {
		envelope = models.GetEnvelope(para.Envelope_id)
	} else {
		for _, value := range d {
			json.Unmarshal([]byte(value), &envelope)
			if envelope.EnvelopeId == commons.ID(para.Envelope_id) {
				break
			}
			// fmt.Println("3 ", commons.GetRedis().LLen(c, "uid"+strconv.Itoa(para.Uid)))
		}
		//移除此元素
		commons.GetRedis().LRem(c, "uid"+strconv.Itoa(para.Uid), 1, envelope)
	}
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
	commons.GetRedis().RPush(c, "uid"+strconv.Itoa(para.Uid), envelope)
	commons.GetRedis().Expire(c, "uid"+strconv.Itoa(para.Uid), 600*1000000000)
	// envelope := models.GetEnvelope(para.Envelope_id)

	//4.TODO 更新红包的状态
	// models.UpdateState(para.Envelope_id, 1)
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
	// envelopes := models.GetEnvelopeList(para.Uid)
	d, err := commons.GetRedis().LRange(c, "uid"+strconv.Itoa(para.Uid), 0, -1).Result()
	// fmt.Println("2 ", commons.GetRedis().LLen(c, "uid"+strconv.Itoa(para.Uid)))
	var envelopes []models.Envelope
	if err != nil {
		envelopes = models.GetEnvelopeList(para.Uid)
	} else {
		for _, value := range d {
			var stem models.Envelope
			json.Unmarshal([]byte(value), &stem)
			envelopes = append(envelopes, stem)
		}
	}
	//4.TODO 更新缓存
	commons.GetRedis().Expire(c, "uid"+strconv.Itoa(para.Uid), 600*1000000000)
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
