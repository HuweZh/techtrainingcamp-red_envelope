package middleware

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"huhusw.com/red_envelope/commons"
	"huhusw.com/red_envelope/models"
)

func WalletMiddle(c *gin.Context) {
	//请求前处理
	fmt.Println("查看红包列表")
	//1.获取请求参数
	para := commons.GetParamter(c)

	d, _ := commons.GetRedis().LRange(c, "uid"+strconv.Itoa(para.Uid), 0, -1).Result()
	// fmt.Println("2 ", commons.GetRedis().LLen(c, "uid"+strconv.Itoa(para.Uid)))
	var envelopes []models.Envelope
	if len(d) == 0 {
		fmt.Println("查询数据库")
		envelopes = models.GetEnvelopeList(para.Uid)
		fmt.Println(len(envelopes))
	} else {
		for _, value := range d {
			var stem models.Envelope
			json.Unmarshal([]byte(value), &stem)
			envelopes = append(envelopes, stem)
		}
	}
	if len(envelopes) == 0 {
		c.Abort()
		commons.R(c, -1, "此用户没有红包", nil)
	}
	//超时时间的单位为微秒，100*1000000000 是100秒
	c.Set("envelopes", envelopes)
	//中间件通信，设置值
	c.Set("name", "我是中间件中的数据")
	//执行请求
	c.Next()
	//中断请求
	// c.Abort()
	//请求后处理
	fmt.Println("查看红包列表后")
}
