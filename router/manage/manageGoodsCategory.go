package manage

import (
	v1 "Graduation/api/v1"
	"Graduation/middleware"

	"github.com/gin-gonic/gin"
)

type ManageGoodsCategoryRouter struct {
}

func (r *ManageGoodsCategoryRouter) ApiManageGoodsCategoryRouter(Router *gin.RouterGroup) {

	goodsCategoryRouter := Router.Group("v1").Use(middleware.AdminJWTAuth())
	superGoodsCategoryRouter := Router.Group("v1").Use(middleware.AdminJWTAuth(), middleware.SuperAdminJWTAuth()) // 验证超级管理员权限
	var goodsCategoryApi = v1.ApiGroupApp.ManageApiGroup.ManageGoodsCategoryApi
	{
		goodsCategoryRouter.GET("categories", goodsCategoryApi.GetCategoryList) // 获取分类数据列表

		// goodsCategoryRouter.PUT("categories", goodsCategoryApi.UpdateCategory)
		// goodsCategoryRouter.GET("categories/:id", goodsCategoryApi.GetCategory)
		// goodsCategoryRouter.DELETE("categories", goodsCategoryApi.DelCategory)
		// goodsCategoryRouter.GET("categories4Select", goodsCategoryApi.ListForSelect)
	}
	// 超级管理员
	{
		superGoodsCategoryRouter.POST("categories", goodsCategoryApi.CreateCategory) // 创建商品种类分类
	}
}
