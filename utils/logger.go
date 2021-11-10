package utils

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"path"
)
var (
	logPath = "./logs"
	logFile = "gin.log"
)
var MyLog = logrus.New()

// InitLogger 日志初始化
func InitLogger()  {
	// 打开文件
	logFileName := path.Join(logPath, logFile)
	// 使用滚动压缩方式记录日志
	rolling(logFileName)
	// 设置日志输出JSON格式
	MyLog.SetFormatter(&logrus.JSONFormatter{})
	//MyLog.SetFormatter(&logrus.TextFormatter{})
	// 设置日志记录级别
	MyLog.SetLevel(logrus.DebugLevel)
}
// 日志滚动设置
func rolling(logFile string)  {
	// 设置输出
	MyLog.SetOutput(&lumberjack.Logger{
		Filename:logFile, //日志文件位置
		MaxSize: 1,// 单文件最大容量,单位是MB
		MaxBackups: 10,// 最大保留过期文件个数
		MaxAge: 3 ,// 保留过期文件的最大时间间隔,单位是天
		Compress: true,// 是否需要压缩滚动日志, 使用的 gzip 压缩
	})
}
