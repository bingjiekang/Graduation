package response

import "time"

// 用户订单返回信息
type MallOrderResponse struct {
	OrderId                int               `json:"orderId"`
	OrderNo                string            `json:"orderNo"`
	TotalPrice             int               `json:"totalPrice"`
	PayType                int               `json:"payType"`
	OrderStatus            int               `json:"orderStatus"`
	OrderStatusString      string            `json:"orderStatusString"`
	CreatedAt              time.Time         `json:"createTime"`
	NewBeeMallOrderItemVOS []MallOrderItemVO `json:"newBeeMallOrderItemVOS"`
}

type MallOrderItemVO struct {
	GoodsId       int    `json:"goodsId"`
	GoodsName     string `json:"goodsName"`
	GoodsCount    int    `json:"goodsCount"`
	GoodsCoverImg string `json:"goodsCoverImg"`
	SellingPrice  int    `json:"sellingPrice"`
}

type MallOrderDetailVO struct {
	OrderNo                string            `json:"orderNo"`
	TotalPrice             int               `json:"totalPrice"`
	PayStatus              int               `json:"payStatus"`
	PayType                int               `json:"payType"`
	PayTypeString          string            `json:"payTypeString"`
	PayTime                time.Time         `json:"payTime"`
	OrderStatus            int               `json:"orderStatus"`
	OrderStatusString      string            `json:"orderStatusString"`
	CreatedAt              time.Time         `json:"createTime"`
	NewBeeMallOrderItemVOS []MallOrderItemVO `json:"newBeeMallOrderItemVOS"`
}
