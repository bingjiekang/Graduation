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
		// 注册并初始化商城路由
		mallRouter.ApiMallUserRouter(MallGroup)          // 注册初始化用户登陆及注册信息路由
		mallRouter.ApiMallUserAddressRouter(MallGroup)   // 注册并初始化用户地址路由
		mallRouter.ApiMallIndexRouter(MallGroup)         // 注册并初始化首页信息路由
		mallRouter.ApiMallGoodsCategoryRouter(MallGroup) // 注册并初始化分类页信息路由
		mallRouter.ApiMallGoodsInfoRouter(MallGroup)     // 注册并初始化商品信息路由
		mallRouter.ApiMallShopCartRouter(MallGroup)      // 注册并初始化购物车信息路由
		mallRouter.ApiMallOrderRouter(MallGroup)         // 注册并初始化订单路由
	}
	return Router
}
