package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitMiddle(c *gin.Context) {
	//检测作弊用户
	//请求前处理
	fmt.Println("该用户未作弊，执行请求")

	//中间件通信，设置值
	c.Set("name", "我是中间件中的数据")
	//执行请求
	c.Next()
	//请求后处理
	fmt.Println("请求执行结束....")
}
