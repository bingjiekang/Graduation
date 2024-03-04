package mall

import "Graduation/service"

type MallGroup struct {
	MallUserApi
	MallUserAddressApi
	MallIndexApi
	MallGoodsCategoryApi
}

var (
	mallUserService          = service.ServiceGroupApp.MallServiceGroup.MallUserService            // 处理用户登录/注册/退出
	mallUserAddressService   = service.ServiceGroupApp.MallServiceGroup.MallUserAddressService     // 处理用户对地址的操作
	mallCarouselService      = service.ServiceGroupApp.MallServiceGroup.MallCarouselService        // 用来获取首页轮播图
	mallIndexConfigService   = service.ServiceGroupApp.MallServiceGroup.MallIndexInfomationService // 获取首页新品上线/热门商品/最新推荐
	mallGoodsCategoryService = service.ServiceGroupApp.MallServiceGroup.MallGoodsCategoryService   // 获取分类页信息
)
