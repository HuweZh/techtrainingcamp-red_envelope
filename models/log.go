package models

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

// 日志记录到文件
func init() {
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	gin.DisableConsoleColor()

	// 记录到文件。
	f, _ := os.Create("./log/logger.log")

	gin.DefaultWriter = io.MultiWriter(f)
}
