package service

import (
	"github.com/gin-gonic/gin"
	"kube-control-api/api/common/constants"
	"kube-control-api/api/workqueue"
	"net/http"
)

/**
会下载一个打包好的docker 镜像到本地，此镜像会用于搭建kubernetes 集群，
集群配置信息：
*/
func Download_kube_resource(c *gin.Context) {

	downloadResourceTask := workqueue.Task{
		Args:    []string{constants.GET_KUBE_DOCKER_IMAGE_ADDRESS(), constants.GET_KUBE_VERSION()},
		Command: "scripts/pull-resource-package.sh",
		Name:    "download resource aaaaaaaaaaaaaa",
		Type:    "download resource",
	}
	workqueue.Evaluator.AddWork(&downloadResourceTask)

	c.JSON(http.StatusOK, gin.H{
		"status": "200",
	})
}
