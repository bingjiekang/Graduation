package manage

import "Graduation/service"

type ManageGroup struct {
	ManageAdminUserApi
	ManageGoodsCategoryApi
}

var (
	manageUserService            = service.ServiceGroupApp.ManageServiceGroup.ManageUserService            // 管理员以及超级管理员操作
	mallGoodsCategoryService     = service.ServiceGroupApp.ManageServiceGroup.ManageGoodsCategoryService   // 分类商品管理
	fileUploadAndDownloadService = service.ServiceGroupApp.ManageServiceGroup.FileUploadAndDownloadService // 上传图片
)
