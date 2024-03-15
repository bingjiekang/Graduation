package mall

import (
	"Graduation/global"
	"Graduation/model/common/response"
	"Graduation/utils/enum"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MallIndexApi struct {
}

// 首页信息转接口
func (m MallIndexApi) MallIndexInfomation(c *gin.Context) {
	// 轮播商品展示
	err, _, mallCarouseInfo := mallCarouselService.GetIndexCarousels(5)
	if err != nil {
		global.GVA_LOG.Error("轮播图获取失败"+err.Error(), zap.Error(err))
		response.FailWithMessage("轮播图获取失败", c)
	}
	// 新品上线展示
	err, newGoodses := mallIndexConfigService.GetIndexInfomation(enum.IndexGoodsNew.Code(), 5)
	if err != nil {
		global.GVA_LOG.Error("新品获取失败"+err.Error(), zap.Error(err))
		response.FailWithMessage("新品获取失败", c)
	}
	// 热门商品展示
	err, hotGoodses := mallIndexConfigService.GetIndexInfomation(enum.IndexGoodsHot.Code(), 4)
	if err != nil {
		global.GVA_LOG.Error("热门商品获取失败"+err.Error(), zap.Error(err))
		response.FailWithMessage("热门商品获取失败", c)
		return
	}
	// 最新推荐商品展示
	err, recommendGoodses := mallIndexConfigService.GetIndexInfomation(enum.IndexGoodsRecommond.Code(), 10)
	if err != nil {
		global.GVA_LOG.Error("推荐商品获取失败"+err.Error(), zap.Error(err))
		response.FailWithMessage("推荐商品获取失败", c)
	}
	// 首页全部商品数据
	indexResult := make(map[string]interface{})
	indexResult["carousels"] = mallCarouseInfo         // 轮播图数据
	indexResult["newGoodses"] = newGoodses             // 新品上市
	indexResult["hotGoodses"] = hotGoodses             // 热门商品
	indexResult["recommendGoodses"] = recommendGoodses // 推荐商品
	response.OkWithData(indexResult, c)
}
