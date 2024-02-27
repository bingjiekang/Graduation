package mall

import (
	v1 "Graduation/api/v1"

	"github.com/gin-gonic/gin"
)

type MallUserRouter struct {
}

// 用户信息登录路由
func (m *MallUserRouter) ApiMallUserRouter(router *gin.RouterGroup) {

	userRoute := router.Group("v1")
	var mallUserApi = v1.ApiGroupApp.MallApiGroup.MallUserApi
	{
		userRoute.POST("/user/register", mallUserApi.UserRegister) // 用户注册
		userRoute.POST("/user/login", mallUserApi.UserLogin)       // 用户登陆
	}
	{

	}

}
