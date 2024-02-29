package mall

import (
	v1 "Graduation/api/v1"
	"Graduation/middleware"

	"github.com/gin-gonic/gin"
)

type MallUserRouter struct {
}

// 用户信息登录路由
func (m *MallUserRouter) ApiMallUserRouter(router *gin.RouterGroup) {
	mallUserRouter := router.Group("v1").Use(middleware.UserJWTAuth())
	userRoute := router.Group("v1")
	var mallUserApi = v1.ApiGroupApp.MallApiGroup.MallUserApi
	{
		userRoute.POST("/user/register", mallUserApi.UserRegister) // 用户注册
		userRoute.POST("/user/login", mallUserApi.UserLogin)       // 用户登陆
	}
	{
		mallUserRouter.POST("/user/logout", mallUserApi.UserLogout)  // 用户登出
		mallUserRouter.GET("/user/info", mallUserApi.UserInfo)       // 用户获取信息
		mallUserRouter.PUT("/user/info", mallUserApi.UpdateUserInfo) // 用户修改信息
	}

}
