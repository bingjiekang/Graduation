package middleware

import (
	"Graduation/global"
	"Graduation/model/common/response"
	"Graduation/service"
	"Graduation/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

var mallUserService = service.ServiceGroupApp.MallServiceGroup.MallUserService

// 中间件验证token
func UserJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			response.UnLogin(nil, c)
			c.Abort()
			return
		}
		err, okm := mallUserService.ExistUserToken(token)
		if err != nil {
			response.UnLogin(nil, c)
			c.Abort()
			return
		}
		if okm == 1 {
			uuid, err, ok := utils.UndoToken(token)
			if err != nil && ok == 0 { // 解码token出现错误
				global.GVA_LOG.Error("解码token出现错误!")
				response.FailWithDetailed(nil, "解码token出现错误", c)
				c.Abort()
				return
			}
			uid := fmt.Sprintf("%d", uuid)
			response.FailWithDetailed(nil, "授权已过期", c)
			global.GVA_REDIS.Del(global.GVA_CTX, uid)
			c.Abort()
			return
		}

		c.Next()
	}

}
