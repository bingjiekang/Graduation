package middleware

import (
	"Graduation/global"
	"Graduation/model/common/response"
	"Graduation/service"
	"Graduation/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	mallUserService             = service.ServiceGroupApp.MallServiceGroup.MallUserService
	manageAdminUserTokenService = service.ServiceGroupApp.ManageServiceGroup.ManageUserService
)

// 用户中间件验证token
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

// 管理员中间件验证
func AdminJWTAuth() gin.HandlerFunc {
	// fmt.Println("11111")
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			response.UnLogin(nil, c)
			c.Abort()
			return
		}
		err, okm := manageAdminUserTokenService.ExistManageToken(token)
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

// 超级管理员验证
func SuperAdminJWTAuth() gin.HandlerFunc {
	// fmt.Println("22222")
	// 管理员验证后部分需要进行超级管理员验证(token已存在)
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		_, ok := manageAdminUserTokenService.IsSuperManageAdmin(token)
		if !ok { // 出现问题,或者不是超级管理员
			response.FailWithDetailed(nil, "您没有此操作权限!", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
