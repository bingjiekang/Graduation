package mall

import (
	"Graduation/global"
	"Graduation/model/mall"
	"Graduation/model/mall/response"
	"Graduation/model/manage"
	"Graduation/model/manage/request"
	"Graduation/utils"
	"Graduation/utils/enum"
	"errors"

	"github.com/jinzhu/copier"

	"time"
)

type MallOrderService struct {
}

// SaveOrder 保存订单
func (m *MallOrderService) SaveOrder(token string, userAddress mall.MallUserAddress, shopCartItems []response.CartItemResponse) (err error, orderNo string) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户"), orderNo
	}
	uuid, _, _ := utils.UndoToken(token)
	var itemIdList []int
	var goodsIds []int
	for _, cartItem := range shopCartItems {
		itemIdList = append(itemIdList, cartItem.CartItemId)
		goodsIds = append(goodsIds, cartItem.GoodsId)
	}
	var mallGoods []manage.MallGoodsInfo
	global.GVA_DB.Where("goods_id in ? ", goodsIds).Find(&mallGoods)
	//检查是否包含已下架商品
	for _, mallGood := range mallGoods {
		if mallGood.GoodsSellStatus != enum.GOODS_UNDER.Code() {
			return errors.New("已下架，无法生成订单"), orderNo
		}
	}
	mallGoodsMap := make(map[int]manage.MallGoodsInfo)
	for _, mallGood := range mallGoods {
		mallGoodsMap[mallGood.GoodsId] = mallGood
	}
	//判断商品库存
	for _, shopCartItemVO := range shopCartItems {
		//查出的商品中不存在购物车中的这条关联商品数据，直接返回错误提醒
		if _, ok := mallGoodsMap[shopCartItemVO.GoodsId]; !ok {
			return errors.New("购物车数据异常！"), orderNo
		}
		if shopCartItemVO.GoodsCount > mallGoodsMap[shopCartItemVO.GoodsId].StockNum {
			return errors.New("库存不足！"), orderNo
		}
	}
	// 删除购物项
	if len(itemIdList) > 0 && len(goodsIds) > 0 {
		if err = global.GVA_DB.Where("cart_item_id in ?", itemIdList).Updates(mall.MallShopCartItem{IsDeleted: 1}).Error; err == nil {
			var stockNumDTOS []request.StockNumDTO
			copier.Copy(&stockNumDTOS, &shopCartItems)
			for _, stockNumDTO := range stockNumDTOS {
				var goodsInfo manage.MallGoodsInfo
				global.GVA_DB.Where("goods_id =?", stockNumDTO.GoodsId).First(&goodsInfo)
				if err = global.GVA_DB.Where("goods_id =? and stock_num>= ? and goods_sell_status = 0", stockNumDTO.GoodsId, stockNumDTO.GoodsCount).Updates(manage.MallGoodsInfo{StockNum: goodsInfo.StockNum - stockNumDTO.GoodsCount}).Error; err != nil {
					return errors.New("库存不足！"), orderNo
				}
			}
			//生成订单号
			orderNo = utils.GenOrderNo()
			// 明天继续 对商家同时添加订单号
			priceTotal := 0
			//保存订单
			var mallOrder manage.MallOrder
			mallOrder.OrderNo = orderNo
			mallOrder.UUid = uuid
			//总价
			for _, mallShopCartItemVO := range shopCartItems {
				priceTotal = priceTotal + mallShopCartItemVO.GoodsCount*mallShopCartItemVO.SellingPrice
			}
			if priceTotal < 1 {
				return errors.New("订单价格异常！"), orderNo
			}
			mallOrder.TotalPrice = priceTotal
			mallOrder.ExtraInfo = ""
			//生成订单项并保存订单项纪录
			if err = global.GVA_DB.Save(&mallOrder).Error; err != nil {
				return errors.New("订单入库失败！"), orderNo
			}
			//生成订单收货地址快照，并保存至数据库
			var mallOrderAddress mall.MallOrderAddress
			copier.Copy(&mallOrderAddress, &userAddress)
			mallOrderAddress.OrderId = mallOrder.OrderId
			//生成所有的订单项快照，并保存至数据库
			var mallOrderItems []manage.MallOrderItem
			for _, mallShoppingCartItemVO := range shopCartItems {
				var mallOrderItem manage.MallOrderItem
				copier.Copy(&mallOrderItem, &mallShoppingCartItemVO)
				mallOrderItem.OrderId = mallOrder.OrderId
				mallOrderItems = append(mallOrderItems, mallOrderItem)
			}
			if err = global.GVA_DB.Save(&mallOrderItems).Error; err != nil {
				return err, orderNo
			}
		}
	}
	return
}

// PaySuccess 支付订单
func (m *MallOrderService) PaySuccess(orderNo string, payType int) (err error) {
	var mallOrder manage.MallOrder
	err = global.GVA_DB.Where("order_no = ? and is_deleted=0 ", orderNo).First(&mallOrder).Error
	if mallOrder != (manage.MallOrder{}) {
		if mallOrder.OrderStatus != 0 {
			return errors.New("订单状态异常！")
		}
		mallOrder.OrderStatus = enum.ORDER_PAID.Code()
		mallOrder.PayType = payType
		mallOrder.PayStatus = 1
		localSH, _ := time.LoadLocation("Asia/Shanghai")
		mallOrder.PayTime = time.Now().In(localSH)
		err = global.GVA_DB.Save(&mallOrder).Error
	}
	return
}

// FinishOrder 完结订单
func (m *MallOrderService) FinishOrder(token string, orderNo string) (err error) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	uuid, _, _ := utils.UndoToken(token)
	var mallOrder manage.MallOrder
	if err = global.GVA_DB.Where("order_no=? and is_deleted = 0", orderNo).First(&mallOrder).Error; err != nil {
		return errors.New("未查询到记录！")
	}
	if mallOrder.UUid != uuid {
		return errors.New("未查询到您的信息,禁止该操作！")
	}
	mallOrder.OrderStatus = enum.ORDER_SUCCESS.Code()
	err = global.GVA_DB.Save(&mallOrder).Error
	return
}

// CancelOrder 关闭订单
func (m *MallOrderService) CancelOrder(token string, orderNo string) (err error) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	uuid, _, _ := utils.UndoToken(token)
	var mallOrder manage.MallOrder
	if err = global.GVA_DB.Where("order_no=? and is_deleted = 0", orderNo).First(&mallOrder).Error; err != nil {
		return errors.New("未查询到记录！")
	}
	if mallOrder.UUid != uuid {
		return errors.New("未查询到您的信息,禁止该操作！")
	}
	if utils.NumsInList(mallOrder.OrderStatus, []int{enum.ORDER_SUCCESS.Code(),
		enum.ORDER_CLOSED_BY_MALLUSER.Code(), enum.ORDER_CLOSED_BY_EXPIRED.Code(), enum.ORDER_CLOSED_BY_JUDGE.Code()}) {
		return errors.New("订单状态异常！")
	}
	mallOrder.OrderStatus = enum.ORDER_CLOSED_BY_MALLUSER.Code()
	err = global.GVA_DB.Save(&mallOrder).Error
	return
}

// GetOrderDetailByOrderNo 获取订单详情
func (m *MallOrderService) GetOrderDetailByOrderNo(token string, orderNo string) (err error, orderDetail response.MallOrderDetailVO) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户"), orderDetail
	}
	uuid, _, _ := utils.UndoToken(token)
	var mallOrder manage.MallOrder
	if err = global.GVA_DB.Where("order_no=? and is_deleted = 0", orderNo).First(&mallOrder).Error; err != nil {
		return errors.New("未查询到记录！"), orderDetail
	}
	if mallOrder.UUid != uuid {
		return errors.New("未查询到您的信息,禁止该操作！"), orderDetail
	}
	var orderItems []manage.MallOrderItem
	err = global.GVA_DB.Where("order_id = ?", mallOrder.OrderId).Find(&orderItems).Error
	if len(orderItems) <= 0 {
		return errors.New("订单项不存在！"), orderDetail
	}

	var mallOrderItemVOS []response.MallOrderItemVO
	copier.Copy(&mallOrderItemVOS, &orderItems)
	copier.Copy(&orderDetail, &mallOrder)
	// 订单状态前端显示为中文
	_, OrderStatusStr := enum.GetMallOrderStatusEnumByStatus(orderDetail.OrderStatus)
	_, payTapStr := enum.GetMallOrderStatusEnumByStatus(orderDetail.PayType)
	orderDetail.OrderStatusString = OrderStatusStr
	orderDetail.PayTypeString = payTapStr
	orderDetail.NewBeeMallOrderItemVOS = mallOrderItemVOS

	return
}

// MallOrderListBySearch 搜索订单
func (m *MallOrderService) MallOrderListBySearch(token string, pageNumber int, status string) (err error, list []response.MallOrderResponse, total int64) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户"), list, total
	}
	uuid, _, _ := utils.UndoToken(token)
	// 根据搜索条件查询
	var mallOrders []manage.MallOrder
	db := global.GVA_DB.Model(&mallOrders)
	if status != "" {
		db.Where("order_status = ?", status)
	}
	err = db.Where("u_uid =? and is_deleted=0 ", uuid).Count(&total).Error
	// 这里前段没有做滚动加载，直接显示全部订单
	// limit := 5
	offset := 5 * (pageNumber - 1)
	err = db.Offset(offset).Order(" order_id desc").Find(&mallOrders).Error
	var orderListVOS []response.MallOrderResponse
	if total > 0 {
		//数据转换 将实体类转成vo
		copier.Copy(&orderListVOS, &mallOrders)
		//设置订单状态中文显示值
		for _, newBeeMallOrderListVO := range orderListVOS {
			_, statusStr := enum.GetMallOrderStatusEnumByStatus(newBeeMallOrderListVO.OrderStatus)
			newBeeMallOrderListVO.OrderStatusString = statusStr
		}
		// 返回订单id
		var orderIds []int
		for _, order := range mallOrders {
			orderIds = append(orderIds, order.OrderId)
		}
		//获取OrderItem
		var orderItems []manage.MallOrderItem
		if len(orderIds) > 0 {
			global.GVA_DB.Where("order_id in ?", orderIds).Find(&orderItems)
			itemByOrderIdMap := make(map[int][]manage.MallOrderItem)
			for _, orderItem := range orderItems {
				itemByOrderIdMap[orderItem.OrderId] = []manage.MallOrderItem{}
			}
			for k, v := range itemByOrderIdMap {
				for _, orderItem := range orderItems {
					if k == orderItem.OrderId {
						v = append(v, orderItem)
					}
					itemByOrderIdMap[k] = v
				}
			}
			//封装每个订单列表对象的订单项数据
			for _, mallOrderListVO := range orderListVOS {
				if _, ok := itemByOrderIdMap[mallOrderListVO.OrderId]; ok {
					orderItemListTemp := itemByOrderIdMap[mallOrderListVO.OrderId]
					var newBeeMallOrderItemVOS []response.MallOrderItemVO
					copier.Copy(&newBeeMallOrderItemVOS, &orderItemListTemp)
					mallOrderListVO.NewBeeMallOrderItemVOS = newBeeMallOrderItemVOS
					_, OrderStatusStr := enum.GetMallOrderStatusEnumByStatus(mallOrderListVO.OrderStatus)
					mallOrderListVO.OrderStatusString = OrderStatusStr
					list = append(list, mallOrderListVO)
				}
			}
		}
	}
	return err, list, total
}
