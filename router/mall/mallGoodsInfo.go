package mall

import (
	v1 "Graduation/api/v1"

	"github.com/gin-gonic/gin"
)

type MallGoodsInfoRouter struct {
}

func (m *MallGoodsInfoRouter) ApiMallGoodsInfoRouter(Router *gin.RouterGroup) {
	mallGoodsRouter := Router.Group("v1")
	var mallGoodsInfoApi = v1.ApiGroupApp.MallApiGroup.MallGoodsInfoApi
	{
		mallGoodsRouter.GET("/goods/detail/:id", mallGoodsInfoApi.GoodsDetail) //商品详情
		mallGoodsRouter.GET("/search", mallGoodsInfoApi.GoodsSearch)           // 商品搜索
	}
}
