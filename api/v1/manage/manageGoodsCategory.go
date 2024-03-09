package manage

import (
	"Graduation/global"
	"Graduation/model/common/response"
	"Graduation/model/manage/request"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ManageGoodsCategoryApi struct {
}

// CreateCategory 新建商品分类
func (g *ManageGoodsCategoryApi) CreateCategory(c *gin.Context) {
	var category request.MallGoodsCategoryReq
	_ = c.ShouldBindJSON(&category)
	if err := mallGoodsCategoryService.AddCategory(category); err != nil {
		global.GVA_LOG.Error("创建分类失败", zap.Error(err))
		response.FailWithMessage("创建分类失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// GetCategoryList 获取商品分类
func (g *ManageGoodsCategoryApi) GetCategoryList(c *gin.Context) {
	var req request.SearchCategoryParams
	_ = c.ShouldBindQuery(&req)
	if err, list, total := mallGoodsCategoryService.SelectCategoryPage(req); err != nil {
		global.GVA_LOG.Error("获取分类商品失败！", zap.Error(err))
		response.FailWithMessage("获取分类商品失败:"+err.Error(), c)
	} else {
		response.OkWithData(response.PageResult{
			List:       list,
			TotalCount: total,
			CurrPage:   req.PageNumber,
			PageSize:   req.PageSize,
			TotalPage:  int(total) / req.PageSize,
		}, c)
	}
}
