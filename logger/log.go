package logger

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

// 日志记录到文件
func init() {
	// 设置为json格式的日志
	Log.Formatter = &logrus.JSONFormatter{}

	// os.Create("./log/systemgin.log")
	file, err := os.OpenFile("./logger/system.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		Log.Println("创建日志文件/打开日志文件失败")
	}
	// 设置log默认文件输出
	Log.Out = file
	gin.SetMode(gin.ReleaseMode)
	// gin框架自己记录的日志也会输出
	gin.DefaultWriter = Log.Out
	// 设置日志级别
	Log.Level = logrus.InfoLevel
}
