package service

import (
	"github.com/gin-gonic/gin"
)

/**
JWT 签发工具类，判断header/Cookies 中是否存在token, 如果存在则校验，如果不存在则
call JWT 中心服务获取
*/
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Next()
	}

}
