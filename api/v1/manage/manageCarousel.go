package manage

import (
	"Graduation/global"
	"Graduation/model/common/request"
	"Graduation/model/common/response"
	requ "Graduation/model/manage/request"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ManageCarouselApi struct {
}

// 创建轮播图
func (m *ManageCarouselApi) CreateCarousel(c *gin.Context) {
	var req requ.MallCarouselAddParam
	_ = c.ShouldBindJSON(&req)
	if err := mallCarouselService.CreateCarousel(req); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败"+err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// 显示轮播图列表
// GetCarouselList 分页获取MallCarousel列表
func (m *ManageCarouselApi) GetCarouselList(c *gin.Context) {
	var pageInfo requ.MallCarouselSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := mallCarouselService.GetCarouselInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取轮播图列表失败!"+err.Error(), zap.Error(err))
		response.FailWithMessage("获取轮播图列表失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:       list,
			TotalCount: total,
			CurrPage:   pageInfo.PageNumber,
			PageSize:   pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// FindMallCarousel 用id查询MallCarousel
func (m *ManageCarouselApi) FindCarousel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err, mallCarousel := mallCarouselService.GetCarousel(id); err != nil {
		global.GVA_LOG.Error("查询指定轮播图id内容失败!"+err.Error(), zap.Error(err))
		response.FailWithMessage("查询指定轮播图id内容失败", c)
	} else {
		response.OkWithData(mallCarousel, c)
	}
}

// 修改轮播图内容
func (m *ManageCarouselApi) UpdateCarousel(c *gin.Context) {
	var req requ.MallCarouselUpdateParam
	_ = c.ShouldBindJSON(&req)
	if err := mallCarouselService.UpdateCarousel(req); err != nil {
		global.GVA_LOG.Error("更新轮播图失败!", zap.Error(err))
		response.FailWithMessage("更新轮播图失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("更新轮播图成功", c)
	}
}

// 删除轮播图内容
func (m *ManageCarouselApi) DeleteCarousel(c *gin.Context) {
	var ids request.IdsReq
	_ = c.ShouldBindJSON(&ids)
	if err := mallCarouselService.DeleteCarousel(ids); err != nil {
		global.GVA_LOG.Error("删除轮播图失败!", zap.Error(err))
		response.FailWithMessage("删除轮播图失败"+err.Error(), c)
	} else {
		response.OkWithMessage("删除轮播图成功", c)
	}
}
