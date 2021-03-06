package middleware

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"huhusw.com/red_envelope/commons"
	"huhusw.com/red_envelope/models"
)

func SnatchMiddle(c *gin.Context) {
	//1.获取请求参数
	para := commons.GetParamter(c)

	//2.根据id获取用户信息
	var user models.User
	userData, err := commons.GetRedis().Get(c, strconv.Itoa(para.Uid)).Result()
	if err != nil {
		user = models.GetUser(para.Uid)
	} else {
		//反序列化
		json.Unmarshal([]byte(userData), &user)
	}

	if user.MaxCount <= user.CurCount {
		//返回数据
		//此用户的红包抢完了
		c.Abort()
		commons.R(c, commons.CODE_OUT_OF_SNATCH_COUNT_ERROR, commons.RUNOUTOF, nil)
	} else {
		//更新用户的cur_count
		user.CurCount += 1
		//中间件通信，设置值
		c.Set("user", user)

		//执行请求
		c.Next()
	}
	//请求后处理
	//将数据传入写数据库的channel
	models.SetMysqlData(commons.UPDATEUSER, user)
	//更新缓存中的数据
	models.SetRedisData(commons.SET, strconv.Itoa(user.UserId), user, 600*1000000000)
}
