package manage

import (
	"Graduation/global"
	"Graduation/model/common/request"
	mage "Graduation/model/manage"
	"Graduation/model/manage/response"
	"Graduation/utils/enum"
	"errors"
	"strconv"

	"github.com/jinzhu/copier"
)

type ManageOrderService struct {
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

// CloseOrder 商家关闭订单
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
		} else {
			return errors.New("订单不能执行关闭操作")
		}
	}
	return
}

// GetMallOrder 根据id获取MallOrder记录
func (m *ManageOrderService) GetMallOrder(id string) (err error, mallOrderDetailVO response.NewBeeMallOrderDetailVO) {
	var mallOrder mage.MallOrder
	if err = global.GVA_DB.Where("order_id = ?", id).First(&mallOrder).Error; err != nil {
		return
	}
	var orderItems []mage.MallOrderItem
	if err = global.GVA_DB.Where("order_id = ?", mallOrder.OrderId).Find(&orderItems).Error; err != nil {
		return
	}
	//获取订单项数据
	if len(orderItems) > 0 {
		var mallOrderItemVOS []response.NewBeeMallOrderItemVO
		copier.Copy(&mallOrderItemVOS, &orderItems)
		copier.Copy(&mallOrderDetailVO, &mallOrder)

		_, OrderStatusStr := enum.GetMallOrderStatusEnumByStatus(mallOrderDetailVO.OrderStatus)
		_, payTapStr := enum.GetMallOrderStatusEnumByStatus(mallOrderDetailVO.PayType)
		mallOrderDetailVO.OrderStatusString = OrderStatusStr
		mallOrderDetailVO.PayTypeString = payTapStr
		mallOrderDetailVO.NewBeeMallOrderItemVOS = mallOrderItemVOS
	}
	return
}

// GetMallOrderInfoList 分页获取MallOrder记录
func (m *ManageOrderService) GetMallOrderInfoList(info request.PageInfo, orderNo string, orderStatus string) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageNumber - 1)
	// 创建db
	db := global.GVA_DB.Model(&mage.MallOrder{})
	if orderNo != "" {
		db.Where("order_no", orderNo)
	}
	// 0.待支付 1.已支付 2.配货完成 3:出库成功 4.交易成功 -1.手动关闭 -2.超时关闭 -3.商家关闭
	if orderStatus != "" {
		status, _ := strconv.Atoi(orderStatus)
		db.Where("order_status", status)
	}
	var mallOrders []mage.MallOrder
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("update_time desc").Find(&mallOrders).Error
	return err, mallOrders, total
}
