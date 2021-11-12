package commons

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"huhusw.com/red_envelope/logger"
)

//url请求中附带的参数
type Paramter struct {
	Uid         int `json:"uid"`
	Envelope_id int `json:"envelope_id"`
}

//从请求中获取到携带的参数
func GetParamter(c *gin.Context) Paramter {
	//接受请求参数
	para := Paramter{}

	//解析参数
	// err := c.ShouldBindBodyWith(&user, binding.JSON)
	err := c.ShouldBindJSON(&para)

	//请求参数错误
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"请求参数出错": err.Error(),
		})
		// fmt.Printf("request paramter error [Err:%s]\n", err.Error())
	}

	//将携带的参数返回
	return para
}
