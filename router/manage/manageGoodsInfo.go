package manage

import (
	v1 "Graduation/api/v1"
	"Graduation/middleware"

	"github.com/gin-gonic/gin"
)

type ManageGoodsInfoRouter struct {
}

func (m *ManageGoodsInfoRouter) ApiManageGoodsInfoRouter(Router *gin.RouterGroup) {
	mallGoodsInfoRouter := Router.Group("v1").Use(middleware.AdminJWTAuth())
	superMallGoodsInfoRouter := Router.Group("v1").Use(middleware.AdminJWTAuth(), middleware.SuperAdminJWTAuth()) // 验证超级管理员权限
	var mallGoodsInfoApi = v1.ApiGroupApp.ManageApiGroup.ManageGoodsInfoApi
	{ // 管理员超级管理员均可
		mallGoodsInfoRouter.POST("goods", mallGoodsInfoApi.CreateGoodsInfo)                    // 新建MallGoodsInfo
		mallGoodsInfoRouter.GET("goods/list", mallGoodsInfoApi.GetGoodsInfoList)               // 获取MallGoodsInfo列表
		mallGoodsInfoRouter.GET("goods/:id", mallGoodsInfoApi.FindGoodsInfo)                   // 根据ID获取MallGoodsInfo
		mallGoodsInfoRouter.PUT("goods", mallGoodsInfoApi.UpdateGoodsInfo)                     // 更新MallGoodsInfo
		mallGoodsInfoRouter.PUT("goods/status/:status", mallGoodsInfoApi.ChangeGoodsInfoByIds) // 上下架
	}
	// 超级管理员
	{
		superMallGoodsInfoRouter.DELETE("deleteMallGoodsInfo", mallGoodsInfoApi.DeleteGoodsInfo) // 删除MallGoodsInfo
	}
}
