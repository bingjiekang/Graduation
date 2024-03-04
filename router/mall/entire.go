package mall

// import mallUser "Graduation/api/v1/mall"

// var MallUser mallUser.MallUser

type MallRouterGroup struct {
	MallUserRouter
	MallUserAddressRouter
	MallIndexRouter
	MallGoodsCategoryRouter
	MallGoodsInfoRouter
	MallShopCartRouter
}
