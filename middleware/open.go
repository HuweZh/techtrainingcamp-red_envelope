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

	//2.根据id获取用户信息
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
	if envelope.Opened == 1 {
		c.Abort()
		commons.R(c, -1, "红包已开", nil)
		commons.GetRedis().RPush(c, "uid"+strconv.Itoa(para.Uid), envelope)
		commons.GetRedis().Expire(c, "uid"+strconv.Itoa(para.Uid), 600*1000000000)
	}
	//超时时间的单位为微秒，100*1000000000 是100秒
	c.Set("envelope", envelope)
	//中间件通信，设置值
	c.Set("name", "我是中间件中的数据")
	//执行请求
	c.Next()
	//中断请求
	// c.Abort()
	//请求后处理
	fmt.Println("开红包之后的判断....")
}
