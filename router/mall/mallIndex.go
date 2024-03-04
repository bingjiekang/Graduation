package mall

import (
	v1 "Graduation/api/v1"

	"github.com/gin-gonic/gin"
)

type MallIndexRouter struct {
}

// 包括上方的轮播图和下方的新品/热门/推荐
func (m *MallIndexRouter) ApiMallIndexRouter(Router *gin.RouterGroup) {
	mallCarouselRouter := Router.Group("v1")
	var mallCarouselApi = v1.ApiGroupApp.MallApiGroup.MallIndexApi
	{
		mallCarouselRouter.GET("index-infos", mallCarouselApi.MallIndexInfomation) // 获取首页数据
	}
}
