package api

import (
	"github.com/gin-gonic/gin"
	"kube-control-api/api/common/constants"
	"kube-control-api/api/service"
)

func SetupRouter() {
	router := gin.Default()
	rootPath := router.Group("/")

	apiPath := rootPath.Group("/api")
	apiPath.POST("/downloadResource", service.Download_kube_resource)

	router.Run(constants.GetEnvOrDefault("KUBOARD_SPRAY_PORT", ":8080"))
}
