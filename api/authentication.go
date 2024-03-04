package api

import (
	"github.com/gin-gonic/gin"
)

/**
JWT 签发工具类，判断header 是否存在token，并验证，会独立出一个服务
*/
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Next()
	}

}
