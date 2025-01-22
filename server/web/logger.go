/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  logger.go  logger.go 2022-11-30
 */

package web

import (
	"github.com/lshaofan/cb-framework/core/logger"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func NewLogger(args interface{}) *Logger {
	l := logrus.New()
	l.SetLevel(logrus.InfoLevel)
	writer, err := logger.GetOutput("http", "exception")
	if err != nil {
		panic(err)
	}
	l.SetOutput(writer)
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return &Logger{logger: l}
}

// AddErrorLog 添加错误日志
func (l *Logger) AddErrorLog(fields map[string]interface{}) {
	l.logger.WithFields(fields).Error()
}

// AddInfoLog 添加信息日志
func (l *Logger) AddInfoLog(fields map[string]interface{}) {
	l.logger.WithFields(fields).Info()
}
