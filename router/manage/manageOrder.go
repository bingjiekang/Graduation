package manage

import (
	v1 "Graduation/api/v1"
	"Graduation/middleware"

	"github.com/gin-gonic/gin"
)

type ManageOrderRouter struct {
}

func (r *ManageOrderRouter) ApiManageOrderRouter(Router *gin.RouterGroup) {
	mallOrderRouter := Router.Group("v1").Use(middleware.AdminJWTAuth())
	superMallOrderRouter := Router.Group("v1").Use(middleware.AdminJWTAuth(), middleware.SuperAdminJWTAuth()) // 验证超级管理员权限
	var mallOrderApi = v1.ApiGroupApp.ManageApiGroup.ManageOrderApi
	// 商家只能查看状态
	{
		mallOrderRouter.GET("orders/:orderId", mallOrderApi.FindMallOrder) // 根据ID获取MallOrder
		mallOrderRouter.GET("orders", mallOrderApi.GetMallOrderList)       // 获取MallOrder列表
	}
	// 商品在超级管理员(总销商)负责发货/出库/取消/
	{
		superMallOrderRouter.PUT("orders/checkDone", mallOrderApi.CheckDoneOrder) // 发货
		superMallOrderRouter.PUT("orders/checkOut", mallOrderApi.CheckOutOrder)   // 出库
		superMallOrderRouter.PUT("orders/close", mallOrderApi.CloseOrder)         // 取消订单
	}
}
