package manage

import (
	"Graduation/global"
	requ "Graduation/model/common/request"
	"Graduation/model/common/response"
	"Graduation/model/manage/request"
	"strconv"

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

// GetCategory 通过id获取分类数据(用来获取选择的数据)
func (g *ManageGoodsCategoryApi) GetCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err, goodsCategory := mallGoodsCategoryService.SelectCategoryById(id)
	if err != nil {
		global.GVA_LOG.Error("获取失败！", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
	} else {
		// response.OkWithData(resp.GoodsCategoryResponse{GoodsCategory: goodsCategory}, c)
		response.OkWithData(goodsCategory, c)
	}
}

// UpdateCategory 修改商品分类信息
func (g *ManageGoodsCategoryApi) UpdateCategory(c *gin.Context) {
	var category request.MallGoodsCategoryReq
	_ = c.ShouldBindJSON(&category)
	if err := mallGoodsCategoryService.UpdateCategory(category); err != nil {
		global.GVA_LOG.Error("更新分类失败", zap.Error(err))
		response.FailWithMessage("更新分类失败，存在相同分类", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// DelCategory 删除(禁用)分类
func (g *ManageGoodsCategoryApi) DelCategory(c *gin.Context) {
	var ids requ.IdsReq
	_ = c.ShouldBindJSON(&ids)
	if err, _ := mallGoodsCategoryService.DeleteGoodsCategoriesByIds(ids); err != nil {
		global.GVA_LOG.Error("删除失败！", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}
