package manage

import "Graduation/service"

type ManageGroup struct {
	ManageAdminUserApi     // 管理员和超级管理员
	ManageGoodsCategoryApi // 分类
	ManageCarouselApi      // 轮播图
	ManageIndexConfigApi   // 首页商品热销/新品/推荐配置
}

var (
	manageUserService            = service.ServiceGroupApp.ManageServiceGroup.ManageUserService            // 管理员以及超级管理员操作
	mallGoodsCategoryService     = service.ServiceGroupApp.ManageServiceGroup.ManageGoodsCategoryService   // 分类商品管理
	fileUploadAndDownloadService = service.ServiceGroupApp.ManageServiceGroup.FileUploadAndDownloadService // 上传图片
	mallCarouselService          = service.ServiceGroupApp.ManageServiceGroup.ManageCarouselService        // 首页轮播图
	mallIndexConfigService       = service.ServiceGroupApp.ManageServiceGroup.ManageIndexConfigService     // 首页商品热销/新品/商品推荐配置
)
