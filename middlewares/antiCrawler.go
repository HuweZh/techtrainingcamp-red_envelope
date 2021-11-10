package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AntiCrawler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//过滤掉不符合的user-agent
		userAgent := c.Request.Header.Get("User-Agent")
		legalUserAgent := strings.Contains(userAgent, "Mozilla")
		if !legalUserAgent {
			c.AbortWithStatus(http.StatusForbidden)
		}

		c.Next()
	}
}
