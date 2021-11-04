package middleware

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"huhusw.com/red_envelope/commons"
	"huhusw.com/red_envelope/models"
)

func SnatchMiddle(c *gin.Context) {
	//请求前处理
	fmt.Println("抢红包之前的判断")
	//1.获取请求参数
	para := commons.GetParamter(c)

	//2.根据id获取用户信息
	var user models.User
	userData, err := commons.GetRedis().Get(c, strconv.Itoa(para.Uid)).Result()
	if err != nil {
		user = models.GetUser(para.Uid)
	} else {
		//更新用户的状态
		json.Unmarshal([]byte(userData), &user)
	}

	//此用户的红包抢完了
	if user.MaxCount <= user.CurCount {
		//返回数据
		c.Abort()
		commons.R(c, commons.BASEERROR, commons.RUNOUTOF, nil)
	} else {
		//更新用户的cur_count
		user.CurCount += 1
	}
	//超时时间的单位为微秒，100*1000000000 是100秒
	c.Set("user", user)
	//中间件通信，设置值
	c.Set("name", "我是中间件中的数据")
	//执行请求
	c.Next()
	//中断请求
	// c.Abort()
	//请求后处理
	fmt.Println("抢红包之后的判断....")

	//将数据传入写数据库的channel
	models.SetMysqlData(commons.UPDATEUSER, user)
	//更新缓存中的数据
	models.SetRedisData(commons.SET, strconv.Itoa(user.UserId), user, 600*1000000000)
}
