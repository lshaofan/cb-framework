package logger

import (
	"fmt"
	//"log"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

//type Option func(*Logger)
//
//func WithFields(fields logrus.Fields) Option {
//	return func(l *Logger) {
//		l.Logger = l.Logger.WithFields(fields).Logger
//	}
//}
//
//func WithPath(path string) Option {
//	return func(l *Logger) {
//		l.path = path
//	}
//}
//
//func WithLevel(level logrus.Level) Option {
//	return func(l *Logger) {
//		l.level = level
//	}
//}
//
//func WithName(name string) Option {
//	return func(l *Logger) {
//		l.name = name
//	}
//}

//type Logger struct {
//	//Log    *logrus.Entry
//	path   string
//	level  logrus.Level
//	name   string
//	Logger *logrus.Logger
//}
//
//func NewLogger(opts ...Option) *Logger {
//	l := &Logger{
//		Logger: logrus.New(),
//	}
//	for _, opt := range opts {
//		opt(l)
//	}
//	// 获取当前时间
//	now := time.Now()
//	// 获取当前年月日
//	year, month, day := now.Date()
//	nowTime := fmt.Sprintf("%d-%d-%d", year, month, day)
//	if l.path == "" {
//		l.path = "info"
//	}
//	l.path = fmt.Sprintf("./runtime/log/%s/%s/%s.log", nowTime, l.path, l.name)
//	lg := logrus.New()
//	lg.SetReportCaller(true)
//	writer, err := rotatelogs.New(
//		l.path+".%Y%m%d",
//		rotatelogs.WithLinkName(l.path),
//		rotatelogs.WithRotationTime(24*time.Hour),  //最小为1分钟轮询。默认60s  低于1分钟就按1分钟来
//		rotatelogs.WithRotationCount(7),            //设置7份 大于7份 或到了清理时间 开始清理
//		rotatelogs.WithRotationSize(100*1024*1024), //设置100MB大小,当大于这个容量时，创建新的日志文件
//
//	)
//	if err != nil {
//		log.Fatalf("Failed to create rotatelogs: %s", err)
//	}
//	lg.SetOutput(writer)
//	lg.SetFormatter(&logrus.JSONFormatter{
//		TimestampFormat: "2006-01-02 15:04:05",
//	})
//	lg.SetLevel(l.level)
//	//l.Log = logrus.NewEntry(lg)
//	l.Logger = lg
//	return l
//}

// Get logrus entry
//func (l *Logger) Get() *logrus.Entry {
//	return l.Logger.
//}

// GetLogger loggers logger
//func (l *Logger) GetLogger() *logrus.Logger {
//	return l.Logger
//}

// GetOutput 获取 SetOutput
func GetOutput(path, name string) (*rotatelogs.RotateLogs, error) {
	// 获取当前时间
	now := time.Now()
	// 获取当前年月日
	year, month, day := now.Date()
	nowTime := fmt.Sprintf("%d-%d-%d", year, month, day)
	filePath := fmt.Sprintf("./runtime/log/%s/%s/%s.log", nowTime, path, name)
	return rotatelogs.New(
		filePath+".%Y%m%d",
		rotatelogs.WithLinkName(filePath),
		rotatelogs.WithRotationTime(24*time.Hour),  //最小为1分钟轮询。默认60s  低于1分钟就按1分钟来
		rotatelogs.WithRotationCount(7),            //设置7份 大于7份 或到了清理时间 开始清理
		rotatelogs.WithRotationSize(100*1024*1024), //设置100MB大小,当大于这个容量时，创建新的日志文件

	)
}
