package mall

import (
	v1 "Graduation/api/v1"
	"Graduation/middleware"

	"github.com/gin-gonic/gin"
)

type MallShopCartRouter struct {
}

func (m *MallUserRouter) ApiMallShopCartRouter(Router *gin.RouterGroup) {
	mallShopCartRouter := Router.Group("v1").Use(middleware.UserJWTAuth())
	var mallShopCartApi = v1.ApiGroupApp.MallApiGroup.MallShopCartApi
	{
		mallShopCartRouter.GET("/shop-cart", mallShopCartApi.CartItemList)                                             // 购物车列表(网页移动端不分页)
		mallShopCartRouter.POST("/shop-cart", mallShopCartApi.AddMallShopCartItem)                                     // 添加购物车商品
		mallShopCartRouter.PUT("/shop-cart", mallShopCartApi.UpdateMallShopCartItem)                                   // 修改购物车商品
		mallShopCartRouter.DELETE("/shop-cart/:newBeeMallShoppingCartItemId", mallShopCartApi.DelMallShoppingCartItem) // 删除购物车商品
		mallShopCartRouter.GET("/shop-cart/settle", mallShopCartApi.ShopTotal)                                         // 根据购物商品id数组查询购物项明细
	}
}
