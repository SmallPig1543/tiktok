package util

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var LogrusObj *logrus.Logger

func InitLog() {
	if LogrusObj != nil {
		src, _ := setOutputFile()
		// 设置输出
		LogrusObj.Out = src
		return
	}
	// 实例化
	logger := logrus.New()
	src, _ := setOutputFile()
	// 设置输出
	logger.Out = src
	// 设置日志级别
	logger.SetLevel(logrus.InfoLevel)
	// 设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	LogrusObj = logger
}

func setOutputFile() (*os.File, error) {

	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/log/"
	}
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(logFilePath, 0777)
		if err != nil {
			return nil, err
		}
	}
	logFileName := time.Now().Format("2006-01-02") + ".log"
	fileName := logFilePath + logFileName
	_, err = os.Stat(fileName)
	if err != nil {
		_, err := os.Create(fileName)
		if err != nil {
			return nil, err
		}
	}
	src, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	return src, nil
}
