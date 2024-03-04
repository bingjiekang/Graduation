package mall

type MallServiceGroup struct {
	MallUserService            // 用户信息对象
	MallUserAddressService     // 用户地址对象
	MallIndexInfomationService // 首页新品/热门商品/推荐商品对象
	MallCarouselService        // 首页轮播图对象
	MallGoodsCategoryService   // 分类页分类信息对象
	MallGoodsInfoService       // 商品详细信息和商品搜索对象
	MallShopCartService        // 购物车及相关操作对象
}
