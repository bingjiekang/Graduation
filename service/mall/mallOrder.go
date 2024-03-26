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
	"fmt"

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
			//生成订单号
			// fmt.Println("开始打印订单号")
			orderNo = utils.GenOrderNo()
			// fmt.Println("订单号", orderNo)
			// 事务开始
			transaction := global.GVA_DB.Begin()
			// 定义字典存储商家id和对应的获得的收益
			var mValue map[int64]int = make(map[int64]int)
			for _, stockNumDTO := range stockNumDTOS {
				// 商品信息
				var goodsInfo manage.MallGoodsInfo
				// 商家订单信息
				var mallAdminOrder manage.MallAdminOrder
				global.GVA_DB.Where("goods_id =?", stockNumDTO.GoodsId).First(&goodsInfo)
				// 更新库存信息(由于updates 不能更新0值,则转用save)
				{
					goodsInfo.StockNum = goodsInfo.StockNum - stockNumDTO.GoodsCount
				}
				if err = global.GVA_DB.Where("goods_id =? and stock_num>= ? and goods_sell_status = 0", stockNumDTO.GoodsId, stockNumDTO.GoodsCount).Save(goodsInfo).Error; err != nil {
					// 事务回滚
					transaction.Rollback()
					return errors.New(fmt.Sprintf("抱歉 %s 商品库存不足！", goodsInfo.GoodsName)), orderNo
				} else {
					// 对商家同时添加订单号,区块哈希等(自己商品的对应订单号)
					{
						mallAdminOrder.OrderNo = orderNo
						mallAdminOrder.Muid = goodsInfo.UUid
						mallAdminOrder.Buid = uuid
						mValue[goodsInfo.UUid] += stockNumDTO.GoodsCount * goodsInfo.SellingPrice // 某一个商品id下的价值计算
					}
					// 使用库存位对商品总量进行求余 对相应商品和库存进行判断 如果未出售则加入改变 已出售则向下寻找
					var seekCount int = 0
					var goodsIndex []int = make([]int, 0)
					var valBlockChain []string = make([]string, 0)
					var valNumber []int64 = make([]int64, 0)
					// 商品库存位
					prevStock := goodsInfo.PrevStock
					for end := 1; end <= goodsInfo.CommodityStocks; end++ {
						// 找到对应数量的库存
						if seekCount == stockNumDTO.GoodsCount {
							break
						}
						prevStock = prevStock%goodsInfo.CommodityStocks + 1
						// 判断此商品是否出售
						var tempMallBlockChain manage.MallBlockChain
						global.GVA_DB.Where("u_uid = ? and commodity = ? and commodity_stocks = ? and is_sale = 0", goodsInfo.UUid, goodsInfo.GoodsId, prevStock).First(&tempMallBlockChain)
						if tempMallBlockChain != (manage.MallBlockChain{}) {
							seekCount++
							goodsIndex = append(goodsIndex, prevStock)
							valBlockChain = append(valBlockChain, tempMallBlockChain.InitBlockHash)
							valNumber = append(valNumber, tempMallBlockChain.Number)
						}
					}
					if len(goodsIndex) != stockNumDTO.GoodsCount {
						return errors.New("对应有效商品数量获取失败" + err.Error()), orderNo
					}
					for k, valIndex := range goodsIndex {
						// 商品区块交易
						var mallManageTrading manage.MallBlockTrading
						// // 商品区块信息
						// var mallManageBlockChain manage.MallBlockChain
						// 添加商品交易区块地址
						mallManageTrading.OrderNo = orderNo
						mallManageTrading.Commodity = goodsInfo.GoodsId
						mallManageTrading.CommodityStocks = valIndex
						mallManageTrading.SellerUid = goodsInfo.UUid
						mallManageTrading.BuyerUid = uuid
						// 当前商品区块哈希
						// if err = global.GVA_DB.Where("u_uid = ? and commodity = ? and commodity_stocks = ? and is_sale = 0", goodsInfo.UUid, goodsInfo.GoodsId, valIndex).First(&mallManageBlockChain).Error; err != nil {
						// 	return errors.New("商品区块信息获取失败" + err.Error()), orderNo
						// }
						// 初始商品哈希
						mallManageTrading.InitBlockHash = valBlockChain[k]
						// 更新区块交易当前交易哈希(mallManageBlockChain.CurrBlockHash)
						newBlock := utils.GenerateNewBlock(
							utils.Block{
								Index:         int64(valNumber[k]),
								CurrBlockHash: valBlockChain[k],
							},
							utils.BlockUserInfo{
								Muid:    goodsInfo.UUid,
								Buid:    uuid,
								GoodsId: goodsInfo.GoodsId,
								Count:   prevStock,
							},
						)
						// 当前商品区块哈希
						mallManageTrading.CurrBlockHash = newBlock.CurrBlockHash
						// 对商品交易信息进行存储
						if err = global.GVA_DB.Create(&mallManageTrading).Error; err != nil {
							transaction.Rollback()
							return errors.New("区块订单入库失败！"), orderNo
						}
					}
					if err = global.GVA_DB.Where("goods_id =?", stockNumDTO.GoodsId).Updates(manage.MallGoodsInfo{PrevStock: prevStock}).Error; err != nil {
						transaction.Rollback()
						return errors.New("商品销售当前位更新失败！"), orderNo
					}
				}
				// 其他状态确定
				{
					mallAdminOrder.TotalPrice = mValue[goodsInfo.UUid]
					mallAdminOrder.ExtraInfo = ""
					localSH, _ := time.LoadLocation("Asia/Shanghai")
					mallAdminOrder.PayTime = time.Now().In(localSH)

				}
				if err = global.GVA_DB.Where("order_no = ? AND m_uid = ?", mallAdminOrder.OrderNo, mallAdminOrder.Muid).Save(&mallAdminOrder).Error; err != nil {
					transaction.Rollback()
					return err, orderNo
				}
			}
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
			localSH, _ := time.LoadLocation("Asia/Shanghai")
			mallOrder.PayTime = time.Now().In(localSH)
			//生成订单项并保存订单项纪录
			if err = global.GVA_DB.Save(&mallOrder).Error; err != nil {
				transaction.Rollback()
				return errors.New("订单入库失败！" + err.Error()), orderNo
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
				// 获取商品当前更新count位后的信息
				var tempMallGoods manage.MallGoodsInfo
				if err = global.GVA_DB.Where("goods_id = ?", mallShoppingCartItemVO.GoodsId).First(&tempMallGoods).Error; err != nil {
					transaction.Rollback()
					return errors.New("商品信息获取失败!"), orderNo
				}
				mallOrderItem.PrevStock = tempMallGoods.PrevStock - mallShoppingCartItemVO.GoodsCount
				mallOrderItems = append(mallOrderItems, mallOrderItem)
			}
			if err = global.GVA_DB.Save(&mallOrderItems).Error; err != nil {
				transaction.Rollback()
				return err, orderNo
			}
			// 事务提交
			transaction.Commit()
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

// FinishOrder 完结订单(需要增加对商品的区块信息更改)
func (m *MallOrderService) FinishOrder(token string, orderNo string) (err error) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	uuid, _, _ := utils.UndoToken(token)
	var mallOrder manage.MallOrder
	if err = global.GVA_DB.Where("order_no = ? and is_deleted = 0", orderNo).First(&mallOrder).Error; err != nil {
		return errors.New("未查询到记录！")
	}
	if mallOrder.UUid != uuid {
		return errors.New("未查询到您的信息,禁止该操作！")
	}
	// 对商品区块进行信息更改orderno->blocktrading->block_chain
	var blockTradings []manage.MallBlockTrading
	if err = global.GVA_DB.Where("order_no = ? and buyer_uid = ?", orderNo, uuid).Find(&blockTradings).Error; err != nil {
		return errors.New("未查询区块商品信息记录！")
	}
	// var mallOrderItems []manage.MallOrderItem
	// if err = global.GVA_DB.Where("order_id = ? ", mallOrder.OrderId).Find(&mallOrderItems).Error; err != nil {
	// 	return errors.New("未查询商品信息记录到记录！")
	// }
	for _, blockTrading := range blockTradings {
		// 获取对象信息
		var blockChain manage.MallBlockChain
		if err = global.GVA_DB.Where("commodity = ? and commodity_stocks = ?", blockTrading.Commodity, blockTrading.CommodityStocks).First(&blockChain).Error; err != nil {
			return errors.New("未查询商品区块信息！")
		}
		// 更新对象block_chain
		if err = global.GVA_DB.Model(&manage.MallBlockChain{}).Where("commodity = ? and commodity_stocks = ?", blockTrading.Commodity, blockTrading.CommodityStocks).Updates(manage.MallBlockChain{
			CurrBlockHash: blockTrading.CurrBlockHash,
			IsSale:        true,
			Number:        blockChain.Number + 1,
		}).Error; err != nil {
			return errors.New("更新商品区块信息失败！")
		}

	}

	mallOrder.OrderStatus = enum.ORDER_SUCCESS.Code()
	err = global.GVA_DB.Save(&mallOrder).Error
	return
}

// CancelOrder 关闭订单(取消订单)复原库存
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
	fmt.Println("订单状态", mallOrder.OrderStatus)
	if utils.NumsInList(mallOrder.OrderStatus, []int{enum.ORDER_SUCCESS.Code(), enum.ORDER_CLOSED_BY_MALLUSER.Code(), enum.ORDER_CLOSED_BY_EXPIRED.Code(), enum.ORDER_CLOSED_BY_JUDGE.Code()}) {
		fmt.Println("状态值", []int{enum.ORDER_SUCCESS.Code(), enum.ORDER_CLOSED_BY_MALLUSER.Code(), enum.ORDER_CLOSED_BY_EXPIRED.Code(), enum.ORDER_CLOSED_BY_JUDGE.Code()})
		return errors.New("订单状态异常！")
	}
	mallOrder.OrderStatus = enum.ORDER_CLOSED_BY_MALLUSER.Code()
	if err = global.GVA_DB.Save(&mallOrder).Error; err != nil {
		return err
	}
	// 复原商品库存 orderno->orderitem->mallgoodsinfo
	var orderItems []manage.MallOrderItem
	if err = global.GVA_DB.Where("order_id = ?", mallOrder.OrderId).Find(&orderItems).Error; err != nil {
		return errors.New("未查询到商品交易记录！")
	}
	for _, orderItem := range orderItems {
		var mallGoodsInfo manage.MallGoodsInfo
		if err = global.GVA_DB.Where("goods_id = ?", orderItem.GoodsId).First(&mallGoodsInfo).Error; err != nil {
			return errors.New("未查询商品信息！")
		}
		// 更新对象 mallGoodsInfo(由于update不能更新0值,则使用save)
		{
			mallGoodsInfo.StockNum = mallGoodsInfo.StockNum + orderItem.GoodsCount
			// mallGoodsInfo.PrevStock = mallGoodsInfo.PrevStock - orderItem.GoodsCount
		}
		if err = global.GVA_DB.Where("goods_id = ?", orderItem.GoodsId).Save(mallGoodsInfo).Error; err != nil {
			return errors.New("更新商品信息失败！")
		}
	}
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
	// copier.Copy(&mallOrderItemVOS, &orderItems)

	// 从区块链交易记录获取信息
	for _, orderItem := range orderItems {
		mallOrderItemVO := response.MallOrderItemVO{
			GoodsId:       orderItem.GoodsId,
			GoodsName:     orderItem.GoodsName,
			GoodsCount:    orderItem.GoodsCount,
			GoodsCoverImg: orderItem.GoodsCoverImg,
			SellingPrice:  orderItem.SellingPrice,
		}
		var blockTradings []manage.MallBlockTrading
		if err = global.GVA_DB.Where("order_no = ? and commodity = ?", orderNo, orderItem.GoodsId).Find(&blockTradings).Error; err != nil {
			return errors.New("未查询到商品区块交易记录！"), orderDetail
		}
		for _, blockTrading := range blockTradings {
			// // 获取对应位商品区块交易哈希
			// var blockTrading manage.MallBlockTrading
			// if err = global.GVA_DB.Where("order_no = ? and commodity = ? and commodity_stocks = ?", mallOrder.OrderNo, orderItem.GoodsId, i).First(&blockTrading).Error; err != nil {
			// 	return errors.New("未查询到商品区块交易记录！"), orderDetail
			// }
			// 把区块交易哈希加入到前端显示信息
			mallOrderItemVO.BlockChainHash = append(mallOrderItemVO.BlockChainHash, blockTrading.CurrBlockHash)
		}
		mallOrderItemVOS = append(mallOrderItemVOS, mallOrderItemVO)
	}
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
