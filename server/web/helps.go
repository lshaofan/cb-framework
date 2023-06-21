package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func GinLoggerFormatter(param gin.LogFormatterParams) string {
	return fmt.Sprintf("[\033[97;42m%s\033[0m] 路由:\033[97;46m%s\033[0m status:\033[97;42m%d\033[0m time:\033[97;44m%s\033[0m 耗时：\033[97;41m%s\033[0m \n",
		param.Method,
		param.Path,
		param.StatusCode,
		param.TimeStamp.Format("2006-01-02 15:04:05"),
		param.Latency,
	)
}
