package service

import (
	"github.com/gin-gonic/gin"
	"kube-control-api/api/common"
	"kube-control-api/api/oscmd"
	"kube-control-api/config/constants"
	"net/http"
	"os"
	"path/filepath"
)

/**
会下载一个打包好的docker 镜像到本地，此镜像会用于搭建kubernetes 集群，
集群配置信息：
*/
func Download_kube_resource(c *gin.Context) {

	versionDir := filepath.Join(constants.GET_CONFIG_RESOURCE_DIR(), constants.GET_KUBE_VERSION())
	_, errExist := os.ReadDir(versionDir)
	if errExist != nil {
		common.HandleError(c, http.StatusConflict, "目录不存在", errExist)
		return
	}

	resultFunc := func(status oscmd.RunningStatus) (string, error) {

		if status.IfSuccess {
			return "download success", nil
		} else {
			return "下载失败，请查看日志", nil
		}
	}
	downloadResourceTask := oscmd.CommandRunnerConfig{
		Type:       "download resource",
		Name:       constants.GET_KUBE_VERSION(),
		ResultFunc: resultFunc,
		Args:       []string{constants.GET_KUBE_DOCKER_IMAGE_ADDRESS()},
		Command:    "scripts/pull-resource-package.sh",
	}

	if err := downloadResourceTask.RunCommand(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to InstallCluster. ", err)
	}
	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data: gin.H{
			"pid": downloadResourceTask.Pid,
		},
	})

}
