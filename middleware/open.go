package middleware

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"huhusw.com/red_envelope/commons"
	"huhusw.com/red_envelope/models"
)

func OpenMiddle(c *gin.Context) {
	//请求前处理
	fmt.Println("开红包之前的判断")
	//1.获取请求参数
	para := commons.GetParamter(c)

	//获取用户信息
	userData, err := commons.GetRedis().Get(c, strconv.Itoa(para.Uid)).Result()

	var envelope models.Envelope
	var user models.User

	//redis中没有存此用户，查询数据库
	if err != nil {
		user = models.GetUser(para.Uid)
		envelope = models.GetEnvelope(para.Uid)
		//将用户存入数据库
		// models.SetRedisData(commons.SET, strconv.Itoa(user.UserId), user, 100*1000000000)
	} else {
		//解析user数据
		json.Unmarshal([]byte(userData), &user)
		//获取红包数据
		d, _ := commons.GetRedis().LRange(c, "uid"+strconv.Itoa(para.Uid), 0, -1).Result()
		for _, value := range d {
			json.Unmarshal([]byte(value), &envelope)
			if envelope.EnvelopeId == commons.ID(para.Envelope_id) {
				break
			}
		}
		//移除此元素
		models.SetRedisData(commons.LREM, "uid"+strconv.Itoa(para.Uid), envelope, 0)
	}

	if envelope.Opened == 1 {
		c.Abort()
		commons.R(c, commons.BASEERROR, commons.OPENED, nil)
	} else {
		//先修改状态，再传入channel，再加入redis缓存
		envelope.Opened = 1
		user.Amount += envelope.Value
	}
	//中间件通信，设置值
	c.Set("envelope", envelope)

	//执行请求
	c.Next()
	//中断请求
	// c.Abort()
	//请求后处理
	fmt.Println("开红包之后的判断....")
	//更新redis
	models.SetRedisData(commons.SET, strconv.Itoa(user.UserId), user, 600*1000000000)
	models.SetRedisData(commons.RPUSH, "uid"+strconv.Itoa(para.Uid), envelope, 0)
	models.SetRedisData(commons.EXPIRE, "uid"+strconv.Itoa(para.Uid), nil, 600*1000000000)
	//将数据传入写数据库的channel
	models.SetMysqlData(commons.UPDATEAMOUNT, user)
	models.SetMysqlData(commons.UPDATEENVELOPESTATE, envelope)
}
