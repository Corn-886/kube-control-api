package main

import (
	"kube-control-api/api"
	"kube-control-api/api/common/log"
)

// 系统启动类
func main() {
	log.InitLogSetting()
	api.SetupRouter()

}
