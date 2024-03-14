package manage

import (
	"Graduation/global"
	requ "Graduation/model/common/request"
	"Graduation/model/common/response"
	"Graduation/model/manage"
	"Graduation/model/manage/request"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ManageGoodsInfoApi struct {
}

// 创建商品
func (m *ManageGoodsInfoApi) CreateGoodsInfo(c *gin.Context) {
	var mallGoodsInfo request.GoodsInfoAddParam
	_ = c.ShouldBindJSON(&mallGoodsInfo)
	token := c.GetHeader("token")
	if err := mallGoodsInfoService.CreateMallGoodsInfo(token, mallGoodsInfo); err != nil {
		global.GVA_LOG.Error("创建商品失败!", zap.Error(err))
		response.FailWithMessage("创建商品失败!"+err.Error(), c)
	} else {
		response.OkWithMessage("创建商品成功", c)
	}
}

// 显示商品列表
// GetMallGoodsInfoList 分页获取MallGoodsInfo列表
func (m *ManageGoodsInfoApi) GetGoodsInfoList(c *gin.Context) {
	var pageInfo request.MallGoodsInfoSearch
	_ = c.ShouldBindQuery(&pageInfo)
	goodsName := c.Query("goodsName")
	goodsSellStatus := c.Query("goodsSellStatus")
	token := c.GetHeader("token")
	if err, list, total := mallGoodsInfoService.GetMallGoodsInfoInfoList(token, pageInfo, goodsName, goodsSellStatus); err != nil {
		global.GVA_LOG.Error("获取商品列表失败!", zap.Error(err))
		response.FailWithMessage("获取商品列表失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:       list,
			TotalCount: total,
			CurrPage:   pageInfo.PageNumber,
			PageSize:   pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// FindMallGoodsInfo 用id和token查询MallGoodsInfo
func (m *ManageGoodsInfoApi) FindGoodsInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	token := c.GetHeader("token")
	err, goodsInfo := mallGoodsInfoService.GetMallGoodsInfo(token, id)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败"+err.Error(), c)
	}
	goodsInfoRes := make(map[string]interface{})
	goodsInfoRes["goods"] = goodsInfo
	if _, thirdCategory := mallGoodsCategoryService.SelectCategoryById(goodsInfo.GoodsCategoryId); thirdCategory != (manage.MallGoodsCategory{}) {
		goodsInfoRes["thirdCategory"] = thirdCategory
		if _, secondCategory := mallGoodsCategoryService.SelectCategoryById(thirdCategory.ParentId); secondCategory != (manage.MallGoodsCategory{}) {
			goodsInfoRes["secondCategory"] = secondCategory
			if _, firstCategory := mallGoodsCategoryService.SelectCategoryById(secondCategory.ParentId); firstCategory != (manage.MallGoodsCategory{}) {
				goodsInfoRes["firstCategory"] = firstCategory
			}
		}
	}
	response.OkWithData(goodsInfoRes, c)

}

// UpdateMallGoodsInfo 更新 MallGoodsInfo
func (m *ManageGoodsInfoApi) UpdateGoodsInfo(c *gin.Context) {
	var mallGoodsInfo request.GoodsInfoUpdateParam
	_ = c.ShouldBindJSON(&mallGoodsInfo)
	token := c.GetHeader("token")
	if err := mallGoodsInfoService.UpdateMallGoodsInfo(token, mallGoodsInfo); err != nil {
		global.GVA_LOG.Error("更新商品信息失败!", zap.Error(err))
		response.FailWithMessage("更新商品信息失败"+err.Error(), c)
	} else {
		response.OkWithMessage("更新商品信息成功", c)
	}
}

// ChangeMallGoodsInfoByIds 修改商品状态 MallGoodsInfo
func (m *ManageGoodsInfoApi) ChangeGoodsInfoByIds(c *gin.Context) {
	var IDS requ.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	sellStatus := c.Param("status")
	token := c.GetHeader("token")
	if err := mallGoodsInfoService.ChangeMallGoodsInfoByIds(token, IDS, sellStatus); err != nil {
		global.GVA_LOG.Error("修改商品状态失败!", zap.Error(err))
		response.FailWithMessage("修改商品状态失败"+err.Error(), c)
	} else {
		response.OkWithMessage("修改商品状态成功", c)
	}
}

// DeleteMallGoodsInfo 超级管理员才能删除商品信息 MallGoodsInfo
func (m *ManageGoodsInfoApi) DeleteGoodsInfo(c *gin.Context) {
	var ids requ.IdsReq
	_ = c.ShouldBindJSON(&ids)
	token := c.GetHeader("token")
	if err, _ := mallGoodsInfoService.DeleteMallGoodsInfo(token, ids); err != nil {
		global.GVA_LOG.Error("删除商品失败!", zap.Error(err))
		response.FailWithMessage("删除商品失败"+err.Error(), c)
	} else {
		response.OkWithMessage("删除商品成功", c)
	}
}
