package mall

import (
	v1 "Graduation/api/v1"

	"github.com/gin-gonic/gin"
)

type MallGoodsCategoryRouter struct {
}

func (m *MallGoodsCategoryRouter) ApiMallGoodsCategoryRouter(Router *gin.RouterGroup) {
	mallGoodsRouter := Router.Group("v1")
	var mallGoodsCategoryApi = v1.ApiGroupApp.MallApiGroup.MallGoodsCategoryApi
	{
		mallGoodsRouter.GET("categories", mallGoodsCategoryApi.GetGoodsCategorize) // 获取分类数据
	}
}
