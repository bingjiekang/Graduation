package manage

import (
	v1 "Graduation/api/v1"
	"Graduation/middleware"

	"github.com/gin-gonic/gin"
)

type ManageCarouselRouter struct {
}

// 首页轮播图相关
func (r *ManageCarouselRouter) ApiManageCarouselRouter(Router *gin.RouterGroup) {
	carouselRouter := Router.Group("v1").Use(middleware.AdminJWTAuth())
	superCarouselRouter := Router.Group("v1").Use(middleware.AdminJWTAuth(), middleware.SuperAdminJWTAuth()) // 验证超级管理员权限
	var mallCarouselApi = v1.ApiGroupApp.ManageApiGroup.ManageCarouselApi
	// 普通管理员
	{
		carouselRouter.GET("carousels", mallCarouselApi.GetCarouselList) // 获取轮播图列表
	}
	// 超级管理员
	{
		superCarouselRouter.POST("carousels", mallCarouselApi.CreateCarousel)   // 新建MallCarousel
		superCarouselRouter.GET("carousels/:id", mallCarouselApi.FindCarousel)  // 根据ID获取轮播图
		superCarouselRouter.PUT("carousels", mallCarouselApi.UpdateCarousel)    // 更新MallCarousel
		superCarouselRouter.DELETE("carousels", mallCarouselApi.DeleteCarousel) // 删除MallCarousel
	}
}
