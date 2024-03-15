package mall

import (
	"Graduation/global"
	"Graduation/model/mall"
	"Graduation/model/mall/request"
	"Graduation/model/mall/response"
	"Graduation/model/manage"
	"Graduation/utils"
	"errors"

	"github.com/jinzhu/copier"
)

type MallShopCartService struct {
}

// 获取购物车信息列表不分页
func (m *MallShopCartService) GetShopCartItems(token string) (err error, cartItems []response.CartItemResponse) {
	var shopCartItems []mall.MallShopCartItem
	var goodsInfos []manage.MallGoodsInfo
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户"), cartItems
	}
	uuid, _, _ := utils.UndoToken(token)
	global.GVA_DB.Where("u_uid=? and is_deleted = 0", uuid).Find(&shopCartItems)
	var goodsIds []int
	for _, shopcartItem := range shopCartItems {
		goodsIds = append(goodsIds, shopcartItem.GoodsId)
	}
	global.GVA_DB.Where("goods_id in ?", goodsIds).Find(&goodsInfos)
	goodsMap := make(map[int]manage.MallGoodsInfo)
	for _, goodsInfo := range goodsInfos {
		goodsMap[goodsInfo.GoodsId] = goodsInfo
	}
	for _, v := range shopCartItems {
		var cartItem response.CartItemResponse
		copier.Copy(&cartItem, &v)
		if _, ok := goodsMap[v.GoodsId]; ok {
			goodsInfo := goodsMap[v.GoodsId]
			cartItem.GoodsName = goodsInfo.GoodsName
			cartItem.GoodsCoverImg = goodsInfo.GoodsCoverImg
			cartItem.SellingPrice = goodsInfo.SellingPrice
		}
		cartItems = append(cartItems, cartItem)
	}

	return err, cartItems
}

// 添加商品到购物车
func (m *MallShopCartService) AddMallCartItem(token string, req request.SaveCartItemParam) (err error) {
	if req.GoodsCount < 1 {
		return errors.New("商品数量不能小于 1 ！")
	}
	if req.GoodsCount > 5 {
		return errors.New("超出单个商品的最大购买数量！")
	}
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	var shopCartItems []mall.MallShopCartItem
	// 是否已存在商品
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid = ? and goods_id = ? and is_deleted = 0", uuid, req.GoodsId).Find(&shopCartItems).Error
	if len(shopCartItems) != 0 || err != nil {
		return errors.New("已存在！无需重复添加！")
	}
	err = global.GVA_DB.Where("goods_id = ? ", req.GoodsId).First(&manage.MallGoodsInfo{}).Error
	if err != nil {
		return errors.New("商品不存在或为空!")
	}
	var total int64
	global.GVA_DB.Where("u_uid = ? and is_deleted = 0", uuid).Count(&total)
	if total > 20 {
		return errors.New("超出购物车最大容量！")
	}
	var shopCartItem mall.MallShopCartItem
	if err = copier.Copy(&shopCartItem, &req); err != nil {
		return err
	}
	shopCartItem.UUid = uuid
	err = global.GVA_DB.Save(&shopCartItem).Error
	return
}

// 更新用户购物车
func (m *MallShopCartService) UpdateMallCartItem(token string, req request.UpdateCartItemParam) (err error) {
	//超出单个商品的最大数量
	if req.GoodsCount > 5 {
		return errors.New("超出单个商品的最大购买数量！")
	}
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	var shopCartItem mall.MallShopCartItem
	if err = global.GVA_DB.Where("cart_item_id=? and is_deleted = 0", req.CartItemId).First(&shopCartItem).Error; err != nil {
		return errors.New("未查询到记录！")
	}
	uuid, _, _ := utils.UndoToken(token)
	if shopCartItem.UUid != uuid {
		return errors.New("未查询到您的信息,禁止该操作！")
	}
	shopCartItem.GoodsCount = req.GoodsCount
	err = global.GVA_DB.Save(&shopCartItem).Error
	return
}

// 删除用户购物车商品
func (m *MallShopCartService) DeleteMallCartItem(token string, id int) (err error) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	var shopCartItem mall.MallShopCartItem
	err = global.GVA_DB.Where("cart_item_id = ? and is_deleted = 0", id).First(&shopCartItem).Error
	if err != nil {
		return
	}
	uuid, _, _ := utils.UndoToken(token)
	if shopCartItem.UUid != uuid {
		return errors.New("未查询到您的信息,禁止该操作！")
	}
	err = global.GVA_DB.Where("cart_item_id = ? and is_deleted = 0", id).UpdateColumns(&mall.MallShopCartItem{IsDeleted: 1}).Error
	return
}

// 获取购物车列表信息总和
func (m *MallShopCartService) GetCartItemsTotal(token string, cartItemIds []int) (err error, cartItemRes []response.CartItemResponse) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户"), cartItemRes
	}
	var shopCartItems []mall.MallShopCartItem
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("cart_item_id in (?) and u_uid = ? and is_deleted = 0", cartItemIds, uuid).Find(&shopCartItems).Error
	if err != nil {
		return
	}
	_, cartItemRes = getMallShopCartItemVOS(shopCartItems)
	//购物车算价
	priceTotal := 0
	for _, cartItem := range cartItemRes {
		priceTotal = priceTotal + cartItem.GoodsCount*cartItem.SellingPrice
	}
	return
}

// 购物车数据转换
func getMallShopCartItemVOS(cartItems []mall.MallShopCartItem) (err error, cartItemsRes []response.CartItemResponse) {
	var goodsIds []int
	for _, cartItem := range cartItems {
		goodsIds = append(goodsIds, cartItem.GoodsId)
	}
	var mallGoods []manage.MallGoodsInfo
	err = global.GVA_DB.Where("goods_id in ?", goodsIds).Find(&mallGoods).Error
	if err != nil {
		return
	}
	mallGoodsMap := make(map[int]manage.MallGoodsInfo)
	for _, goodsInfo := range mallGoods {
		mallGoodsMap[goodsInfo.GoodsId] = goodsInfo
	}
	for _, cartItem := range cartItems {
		var cartItemRes response.CartItemResponse
		copier.Copy(&cartItemRes, &cartItem)
		// 是否包含key
		if _, ok := mallGoodsMap[cartItemRes.GoodsId]; ok {
			mallGoodsTemp := mallGoodsMap[cartItemRes.GoodsId]
			cartItemRes.GoodsCoverImg = mallGoodsTemp.GoodsCoverImg
			goodsName := utils.ReplaceLength(mallGoodsTemp.GoodsName, 28)
			cartItemRes.GoodsName = goodsName
			cartItemRes.SellingPrice = mallGoodsTemp.SellingPrice
			cartItemsRes = append(cartItemsRes, cartItemRes)
		}
	}
	return
}
