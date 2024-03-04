package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	rootPath := router.Group("/")

	apiPath := rootPath.Group("/api", JWTAuthMiddleware())

	apiPath.POST("")

	return router
}
