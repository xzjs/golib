package golib

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Log 封装的log
type Log struct {
	Logger *logrus.Logger
	Fields logrus.Fields
}

// GetLog 获取log
func GetLog(c *gin.Context) Log {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{}

	return Log{
		Logger: log,
		Fields: logrus.Fields{
			"ip":     c.ClientIP(),
			"cookie": c.Request.Cookies(),
			"url":    c.Request.URL.String(),
			"level":  log.Level,
			"time":   time.Now(),
		},
	}
}

// Println 打印日志
func (l *Log) Println(logInfo map[string]interface{}) {
	for k, v := range logInfo {
		l.Fields[k] = v
	}
	log, _ := json.Marshal(l.Fields)
	if err := Producer("log", string(log)); err != nil {
		return
	}
}
