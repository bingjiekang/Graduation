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

type ManageIndexConfigApi struct {
}

// 创建对应热销/新品/商品推荐配置
func (m *ManageIndexConfigApi) CreateIndexConfig(c *gin.Context) {
	var req request.MallIndexConfigAddParams
	_ = c.ShouldBindJSON(&req)
	if err := mallIndexConfigService.CreateMallIndexConfig(req); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败"+err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// 获取对应热销/新品/商品推荐 列表信息
func (m *ManageIndexConfigApi) GetIndexConfigList(c *gin.Context) {
	var pageInfo request.MallIndexConfigSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := mallIndexConfigService.GetMallIndexConfigInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!"+err.Error(), zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:       list,
			TotalCount: total,
			CurrPage:   pageInfo.PageNumber,
			PageSize:   pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// 根据id获取对应信息
func (m *ManageIndexConfigApi) FindIndexConfig(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err, mallIndexConfig := mallIndexConfigService.GetMallIndexConfig(uint(id)); err != nil {
		global.GVA_LOG.Error("查询失败!"+err.Error(), zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(mallIndexConfig, c)
	}
}

// 更新对应商品配置信息
func (m *ManageIndexConfigApi) UpdateIndexConfig(c *gin.Context) {
	var req request.MallIndexConfigUpdateParams
	_ = c.ShouldBindJSON(&req)
	if err := mallIndexConfigService.UpdateMallIndexConfig(req); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// 删除对应商品信息
func (m *ManageIndexConfigApi) DeleteIndexConfig(c *gin.Context) {
	var ids requ.IdsReq
	_ = c.ShouldBindJSON(&ids)
	if err := mallIndexConfigService.DeleteMallIndexConfig(ids); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}
