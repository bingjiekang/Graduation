package manage

import (
	v1 "Graduation/api/v1"
	"Graduation/middleware"

	"github.com/gin-gonic/gin"
)

type ManageIndexConfigRouter struct {
}

func (r *ManageIndexConfigRouter) ApiManageIndexConfigRouter(Router *gin.RouterGroup) {
	indexConfigRouter := Router.Group("v1").Use(middleware.AdminJWTAuth())
	superIndexConfigRouter := Router.Group("v1").Use(middleware.AdminJWTAuth(), middleware.SuperAdminJWTAuth()) // 验证超级管理员权限

	var mallIndexConfigApi = v1.ApiGroupApp.ManageApiGroup.ManageIndexConfigApi
	// 普通管理员
	{
		indexConfigRouter.GET("indexConfigs", mallIndexConfigApi.GetIndexConfigList) // 获取MallIndexConfig列表
	}
	// 超级管理员
	{
		superIndexConfigRouter.POST("indexConfigs", mallIndexConfigApi.CreateIndexConfig)   // 新建MallIndexConfig
		superIndexConfigRouter.GET("indexConfigs/:id", mallIndexConfigApi.FindIndexConfig)  // 根据ID获取MallIndexConfig
		superIndexConfigRouter.PUT("indexConfigs", mallIndexConfigApi.UpdateIndexConfig)    // 更新MallIndexConfig
		superIndexConfigRouter.DELETE("indexConfigs", mallIndexConfigApi.DeleteIndexConfig) // 删除MallIndexConfig
	}

}
