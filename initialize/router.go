package initialize

import (
	"Graduation/middleware"
	"Graduation/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	var Router = gin.Default()

	// 如果有跨域则打开
	Router.Use(middleware.CrossDomain())

	// 前端商城路由
	mallRouter := router.RouterGroupApp.Mall
	// 分组
	MallGroup := Router.Group("api")
	{
		// 初始化商城路由
		mallRouter.ApiMallUserRouter(MallGroup) // 连接到用户登陆及注册信息路由
	}
	return Router
}
