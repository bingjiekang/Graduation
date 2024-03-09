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
		// goodsCategoryRouter.GET("categories4Select", goodsCategoryApi.ListForSelect)
	}
	// 超级管理员
	{
		superGoodsCategoryRouter.POST("categories", goodsCategoryApi.CreateCategory) // 创建商品种类分类
		superGoodsCategoryRouter.GET("categories/:id", goodsCategoryApi.GetCategory) // 获取单个商品登记分类信息
		superGoodsCategoryRouter.PUT("categories", goodsCategoryApi.UpdateCategory)  // 修改商品等级分类信息
		superGoodsCategoryRouter.DELETE("categories", goodsCategoryApi.DelCategory)  // 删除(禁用)商品种类信息
	}
}
