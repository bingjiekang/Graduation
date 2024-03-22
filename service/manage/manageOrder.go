package manage

import (
	"Graduation/global"
	"Graduation/model/common/request"
	"Graduation/model/manage"
	mage "Graduation/model/manage"
	"Graduation/model/manage/response"
	"Graduation/utils"
	"Graduation/utils/enum"
	"errors"
	"strconv"

	"github.com/jinzhu/copier"
)

type ManageOrderService struct {
}

// GetMallOrderInfoList 分页获取 MallOrder 商家记录
func (m *ManageOrderService) GetMallOrderInfoList(token string, info request.PageInfo, orderNo string, orderStatus string) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageNumber - 1)
	// 根据 token 获取对应管理员信息
	uuid, _, _ := utils.UndoToken(token)
	var mallAdminOrder manage.MallAdminUser
	err = global.GVA_DB.Where("u_uid =?", uuid).First(&mallAdminOrder).Error
	if err != nil {
		return errors.New("不存在的管理员用户"), list, total
	}
	// 创建db
	db := global.GVA_DB.Model(&mage.MallOrder{})
	// 判断是否是管理员
	if mallAdminOrder.IsSuperAdmin == 1 {
		if orderNo != "" {
			db.Where("order_no", orderNo)
		}

	} else {
		// 根据 uid 去mall_admin_order 里查询交易订单号
		var mallAdminOrders []manage.MallAdminOrder
		if err = global.GVA_DB.Where("m_uid = ?", uuid).Find(&mallAdminOrders).Error; err != nil {
			return err, list, total
		}
		var orderList []string = make([]string, 0)
		// 将订单号加入到orderlist方便后续查找
		for _, val := range mallAdminOrders {
			orderList = append(orderList, val.OrderNo)
		}
		// 订单号在对应商家订单号表里
		db.Where("order_no in ?", orderList)
		// 传入订单号 查询指定信息
		if orderNo != "" {
			var isOk bool = false
			for _, v := range orderList {
				if v == orderNo {
					isOk = true
					return
				}
			}
			// 没找到对应信息
			if !isOk {
				return errors.New("不可查的商品交易单号,请检查!"), list, total
			}
			db.Where("order_no", orderNo)
		}
	}
	// 0.待支付 1.已支付 2.配货完成 3:出库成功 4.交易成功 -1.手动关闭 -2.超时关闭 -3.商家关闭
	if orderStatus != "" {
		status, _ := strconv.Atoi(orderStatus)
		db.Where("order_status", status)
	}
	var mallOrders []mage.MallOrder
	// 如果有条件搜索 下方会自动创建搜索语句
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&mallOrders).Error
	// 判断超级管理员和商家 返回不同信息
	if mallAdminOrder.IsSuperAdmin != 1 {
		// 对商品价格进行更正
		// 根据登录管理员的信息返回其对应的商品信息
		var mallGoodsInfos []mage.MallGoodsInfo
		if err = global.GVA_DB.Where("u_uid=?", uuid).Find(&mallGoodsInfos).Error; err != nil {
			return
		}
		// 获取对应商家的商品信息
		var tempInfoList []int = make([]int, 0)
		for _, mallGoodsInfo := range mallGoodsInfos {
			tempInfoList = append(tempInfoList, mallGoodsInfo.GoodsId)
		}
		for k, v := range mallOrders {
			var orderItems []mage.MallOrderItem
			if err = global.GVA_DB.Where("order_id = ? and goods_id in ?", v.OrderId, tempInfoList).Find(&orderItems).Error; err != nil {
				return
			}
			var price int
			for _, val := range orderItems {
				price += (val.GoodsCount * val.SellingPrice)
			}
			mallOrders[k].TotalPrice = price
		}
	}
	return err, mallOrders, total
}

// GetMallOrder 根据id获取MallOrder记录
func (m *ManageOrderService) GetMallOrder(token string, id string) (err error, mallOrderDetailVO response.NewBeeMallOrderDetailVO) {
	var mallOrder mage.MallOrder
	if err = global.GVA_DB.Where("order_id = ?", id).First(&mallOrder).Error; err != nil {
		return
	}
	// 检查管理员是否合法
	uuid, _, _ := utils.UndoToken(token)
	var mallAdminOrder manage.MallAdminUser
	err = global.GVA_DB.Where("u_uid =?", uuid).First(&mallAdminOrder).Error
	if err != nil {
		return errors.New("不存在的管理员用户"), mallOrderDetailVO
	}
	var orderItems []mage.MallOrderItem
	if mallAdminOrder.IsSuperAdmin == 1 { // 超级管理员返回全部信息
		if err = global.GVA_DB.Where("order_id = ?", mallOrder.OrderId).Find(&orderItems).Error; err != nil {
			return
		}
	} else {
		// 根据登录管理员的信息返回其对应的商品信息
		var mallGoodsInfos []mage.MallGoodsInfo
		if err = global.GVA_DB.Where("u_uid=?", uuid).Find(&mallGoodsInfos).Error; err != nil {
			return
		}
		// 获取对应商家的商品信息
		var tempInfoList []int = make([]int, 0)
		for _, mallGoodsInfo := range mallGoodsInfos {
			tempInfoList = append(tempInfoList, mallGoodsInfo.GoodsId)
		}
		if err = global.GVA_DB.Where("order_id = ? and goods_id in ?", mallOrder.OrderId, tempInfoList).Find(&orderItems).Error; err != nil {
			return
		}
	}

	//获取订单项数据
	if len(orderItems) > 0 {
		var mallOrderItemVOS []response.NewBeeMallOrderItemVO
		// copier.Copy(&mallOrderItemVOS, &orderItems)
		// 遍历交易订单将order的哈希交易加入到hashblockchain
		for _, orderItem := range orderItems {
			var nums []int = make([]int, 0)
			for count := orderItem.PrevStock + 1; count <= orderItem.PrevStock+orderItem.GoodsCount; count++ {
				// 记录对应的库存位信息
				nums = append(nums, count)
			}
			var tempBlcokTrading []mage.MallBlockTrading
			// 查询区块交易的对应信息
			if err = global.GVA_DB.Where("order_no = ? and Commodity = ? and Commodity_stocks in ?", mallOrder.OrderNo, orderItem.GoodsId, nums).Find(&tempBlcokTrading).Error; err != nil {
				return
			}
			// 将对应信息以及区块哈希交易加入到返回列表里
			var mallOrderItemVO response.NewBeeMallOrderItemVO
			{
				mallOrderItemVO.GoodsId = orderItem.GoodsId
				mallOrderItemVO.GoodsName = orderItem.GoodsName
				mallOrderItemVO.GoodsCount = orderItem.GoodsCount
				mallOrderItemVO.GoodsCoverImg = orderItem.GoodsCoverImg
				mallOrderItemVO.SellingPrice = orderItem.SellingPrice
				for _, val := range tempBlcokTrading {
					// 将当前交易哈希加入
					mallOrderItemVO.HashBlockChain = append(mallOrderItemVO.HashBlockChain, val.CurrBlockHash)
				}
			}
			// 将整体信息返回
			mallOrderItemVOS = append(mallOrderItemVOS, mallOrderItemVO)
		}
		copier.Copy(&mallOrderDetailVO, &mallOrder)
		_, OrderStatusStr := enum.GetMallOrderStatusEnumByStatus(mallOrderDetailVO.OrderStatus)
		_, payTapStr := enum.GetMallOrderStatusEnumByStatus(mallOrderDetailVO.PayType)
		mallOrderDetailVO.OrderStatusString = OrderStatusStr
		mallOrderDetailVO.PayTypeString = payTapStr
		mallOrderDetailVO.NewBeeMallOrderItemVOS = mallOrderItemVOS
	}
	return
}

// CheckDone 修改订单状态为配货成功
func (m *ManageOrderService) CheckDone(ids request.IdsReq) (err error) {
	var orders []mage.MallOrder
	err = global.GVA_DB.Where("order_id in ?", ids.Ids).Find(&orders).Error
	var errorOrders string
	if len(orders) != 0 {
		for _, order := range orders {
			if order.IsDeleted == 1 {
				errorOrders = order.OrderNo + " "
				continue
			}
			if order.OrderStatus != enum.ORDER_PAID.Code() {
				errorOrders = order.OrderNo + " "
			}
		}
		if errorOrders == "" {
			if err = global.GVA_DB.Where("order_id in ?", ids.Ids).UpdateColumns(mage.MallOrder{OrderStatus: 2}).Error; err != nil {
				return err
			}
		} else {
			return errors.New("订单的状态不是支付成功无法执行出库操作")
		}
	}
	return
}

// CheckOut 出库
func (m *ManageOrderService) CheckOut(ids request.IdsReq) (err error) {
	var orders []mage.MallOrder
	err = global.GVA_DB.Where("order_id in ?", ids.Ids).Find(&orders).Error
	var errorOrders string
	if len(orders) != 0 {
		for _, order := range orders {
			if order.IsDeleted == 1 {
				errorOrders = order.OrderNo + " "
				continue
			}
			if order.OrderStatus != enum.ORDER_PAID.Code() && order.OrderStatus != enum.ORDER_PACKAGED.Code() {
				errorOrders = order.OrderNo + " "
			}
		}
		if errorOrders == "" {
			if err = global.GVA_DB.Where("order_id in ?", ids.Ids).UpdateColumns(mage.MallOrder{OrderStatus: 3}).Error; err != nil {
				return err
			}
		} else {
			return errors.New("订单的状态不是支付成功或配货完成无法执行出库操作")
		}
	}
	return
}

// CloseOrder 商家关闭订单()
func (m *ManageOrderService) CloseOrder(ids request.IdsReq) (err error) {
	var orders []mage.MallOrder
	err = global.GVA_DB.Where("order_id in ?", ids.Ids).Find(&orders).Error
	var errorOrders string
	if len(orders) != 0 {
		for _, order := range orders {
			if order.IsDeleted == 1 {
				errorOrders = order.OrderNo + " "
				continue
			}
			if order.OrderStatus == enum.ORDER_SUCCESS.Code() || order.OrderStatus < 0 {
				errorOrders = order.OrderNo + " "
			}
		}
		if errorOrders == "" {
			if err = global.GVA_DB.Where("order_id in ?", ids.Ids).UpdateColumns(mage.MallOrder{OrderStatus: enum.ORDER_CLOSED_BY_JUDGE.Code()}).Error; err != nil {
				return err
			}
			for _, v := range ids.Ids {
				// 复原商品库存 orderno->orderitem->mallgoodsinfo
				var orderItems []manage.MallOrderItem
				if err = global.GVA_DB.Where("order_id = ?", v).Find(&orderItems).Error; err != nil {
					return errors.New("未查询到商品交易记录！")
				}
				for _, orderItem := range orderItems {
					var mallGoodsInfo manage.MallGoodsInfo
					if err = global.GVA_DB.Where("goods_id = ?", orderItem.GoodsId).First(&mallGoodsInfo).Error; err != nil {
						return errors.New("未查询商品信息！")
					}
					// 更新对象 mallGoodsInfo
					if err = global.GVA_DB.Model(&manage.MallGoodsInfo{}).Where("goods_id = ?", orderItem.GoodsId).Updates(manage.MallGoodsInfo{
						StockNum:  mallGoodsInfo.StockNum + orderItem.GoodsCount,
						PrevStock: mallGoodsInfo.PrevStock - orderItem.PrevStock,
					}).Error; err != nil {
						return errors.New("更新商品信息失败！")
					}
				}
			}
		} else {
			return errors.New("订单不能执行关闭操作")
		}
	}
	return
}
