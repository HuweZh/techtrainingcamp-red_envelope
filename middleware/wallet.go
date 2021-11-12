package middleware

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"huhusw.com/red_envelope/commons"
	"huhusw.com/red_envelope/models"
)

func WalletMiddle(c *gin.Context) {
	//获取请求参数
	para := commons.GetParamter(c)
	//获取用户信息
	userData, err := commons.GetRedis().Get(c, strconv.Itoa(para.Uid)).Result()

	var envelopes []models.Envelope
	var user models.User
	//redis中没有存此用户，查询数据库
	if err != nil {
		user = models.GetUser(para.Uid)
		envelopes = models.GetEnvelopeList(para.Uid)
		//存入redis
		for _, envelope := range envelopes {
			models.SetRedisData(commons.RPUSH, "uid"+strconv.Itoa(user.UserId), envelope, 0)
			// commons.GetRedis().RPush(c, "uid"+strconv.Itoa(user.UserId), envelope)
		}
	} else {
		//解析user数据
		json.Unmarshal([]byte(userData), &user)
		//获取红包数据
		d, _ := commons.GetRedis().LRange(c, "uid"+strconv.Itoa(para.Uid), 0, -1).Result()
		for _, value := range d {
			var stem models.Envelope
			json.Unmarshal([]byte(value), &stem)
			envelopes = append(envelopes, stem)
		}
	}

	if len(envelopes) == 0 {
		c.Abort()
		//此用户没有红包
		commons.R(c, commons.NOENVELOPE, commons.HAVEZERO, nil)
	}

	//中间件通信，设置值
	c.Set("envelopes", envelopes)
	c.Set("amount", user.Amount)
	//执行请求
	c.Next()

	//请求后处理
	models.SetRedisData(commons.SET, strconv.Itoa(user.UserId), user, 600*1000000000)
	models.SetRedisData(commons.EXPIRE, "uid"+strconv.Itoa(para.Uid), nil, 600*1000000000)
}
