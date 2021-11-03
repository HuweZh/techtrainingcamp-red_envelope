package commons

import (
	"fmt"

	"github.com/gin-gonic/gin"
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
