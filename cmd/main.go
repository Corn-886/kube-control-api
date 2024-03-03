package main

import (
	"github.com/sirupsen/logrus"
	"log"
)

// 系统启动类
func main() {
	initLogSetting()
}

func initLogSetting() {
	logrus.SetFormatter(new(log.))
}
