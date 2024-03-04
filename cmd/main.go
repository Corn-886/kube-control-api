package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"kube-control-api/api"
	"kube-control-api/config/constants"
	"kube-control-api/config/log"
	"os"
)

// 系统启动类
func main() {
	initLogSetting()

	router := api.SetupRouter()

	port := ":8001"

	router.Run(constants.GetEnvOrDefault("KUBOARD_SPRAY_PORT", port))
}

func initLogSetting() {
	logrus.SetFormatter(new(log.LogrusFormatter)) // 自定义格式输出
	logrus.SetReportCaller(true)                  // 设置日志记录器包含调用者信息
	logrus.SetOutput(os.Stdout)                   // 设置日志输出到控制台

	logLevel := "debug"

	level, err := logrus.ParseLevel(constants.GetEnvOrDefault("KUBOARD_SPRAY_LOGRUS_LEVEL", logLevel))
	if err == nil {
		fmt.Println("设置日志级别为 " + logLevel)
		logrus.SetLevel(level)
	} else {
		fmt.Println("请检查 KUBOARD_SPRAY_LOGRUS_LEVEL 的值，可选的有 panic / fatal / error / warn / info / debug / trace ，当前为： " + logLevel)
		logrus.SetLevel(logrus.InfoLevel)
	}
}
