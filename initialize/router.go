package initialize

import (
	"Graduation/global"
	"Graduation/middleware"
	"Graduation/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	var Router = gin.Default()
	// 为用户头像和文件图片提供静态地址
	Router.StaticFS(global.GVA_CONFIG.Local.Path, http.Dir(global.GVA_CONFIG.Local.Path))
	// 如果有跨域则打开
	Router.Use(middleware.CrossDomain())

	// 商城路由
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
	// 后台管理系统路由
	manageRouter := router.RouterGroupApp.Manage
	// 分组
	ManageGroup := Router.Group("manage-api")
	{
		manageRouter.ApiManageAdminUserRouter(ManageGroup)     // 管理员和超级管理员操作
		manageRouter.ApiManageGoodsCategoryRouter(ManageGroup) // 商品分类
		manageRouter.ApiManageCarouselRouter(ManageGroup)      // 轮播图路由获取
		manageRouter.ApiManageIndexConfigRouter(ManageGroup)   // 商品首页热销/新品/商品推荐配置
		manageRouter.ApiManageGoodsInfoRouter(ManageGroup)     // 商品信息
	}
	return Router
}
