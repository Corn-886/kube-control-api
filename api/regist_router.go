package api

import (
	"github.com/gin-gonic/gin"
	"kube-control-api/api/service"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	rootPath := router.Group("/")

	apiPath := rootPath.Group("/api", service.JWTAuthMiddleware())
	apiPath.POST("/downloadResource", service.Download_kube_resource)

	return router
}
