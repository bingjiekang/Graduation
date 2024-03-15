package mall

import (
	"Graduation/global"
	"Graduation/model/common/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MallGoodsCategoryApi struct {
}

// 返回分类数据 (分类页调用)
func (m *MallGoodsCategoryApi) GetGoodsCategorize(c *gin.Context) {
	err, list := mallGoodsCategoryService.GetGoodsCategories()
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败"+err.Error(), c)
	} else {
		response.OkWithData(list, c)
	}

}
