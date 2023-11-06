package core

import (
	"io"
	"os"
	"strings"
	"vkc/conf"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Hostname string
)

func init() {
	hname, _ := os.Hostname()

	loffset := strings.LastIndex(hname, "-")
	if loffset > 0 {
		Hostname = hname[loffset+1:]
	} else {
		Hostname = hname
	}

	// 设置日志格式
	logrus.SetFormatter(&logrus.TextFormatter{})
	// 设置日志级别
	logrus.SetLevel(logrus.InfoLevel)
	if conf.C.Middleware.LoggerFile == "/dev/null" {
		// 丢弃所有日志
		logrus.SetOutput(io.Discard)
	} else if conf.C.Middleware.LoggerFile != "" {
		// 创建一个新的lumberjack logger
		log := &lumberjack.Logger{
			Filename:   conf.C.Middleware.LoggerFile + hname + ".log",
			MaxSize:    5, // MB
			MaxBackups: 10,
			MaxAge:     14, // days
			Compress:   false,
		}
		logrus.SetOutput(log)
	} else {
		// 默认使用控制台输出
		logrus.SetOutput(os.Stdout)
	}

}
