package mall

import (
	v1 "Graduation/api/v1"
	"Graduation/middleware"

	"github.com/gin-gonic/gin"
)

type MallOrderRouter struct {
}

func (m *MallOrderRouter) ApiMallOrderRouter(Router *gin.RouterGroup) {
	mallOrderRouter := Router.Group("v1").Use(middleware.UserJWTAuth())
	var mallOrderRouterApi = v1.ApiGroupApp.MallApiGroup.MallOrderApi
	{
		mallOrderRouter.GET("/paySuccess", mallOrderRouterApi.PaySuccess)             //模拟支付成功回调的接口
		mallOrderRouter.PUT("/order/:orderNo/finish", mallOrderRouterApi.FinishOrder) //确认收货接口
		mallOrderRouter.PUT("/order/:orderNo/cancel", mallOrderRouterApi.CancelOrder) //取消订单接口
		mallOrderRouter.GET("/order/:orderNo", mallOrderRouterApi.OrderDetailPage)    //订单详情接口
		mallOrderRouter.GET("/order", mallOrderRouterApi.OrderList)                   //订单列表接口
		mallOrderRouter.POST("/saveOrder", mallOrderRouterApi.SaveOrder)
	}
}
