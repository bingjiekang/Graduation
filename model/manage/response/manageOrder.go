package response

import "time"

type NewBeeMallOrderDetailVO struct {
	OrderId                int                     `json:"orderId"`
	OrderNo                string                  `json:"orderNo"`
	TotalPrice             int                     `json:"totalPrice"`
	PayType                int                     `json:"payType"`
	PayTypeString          string                  `json:"payTypeString"`
	OrderStatus            int                     `json:"orderStatus"`
	OrderStatusString      string                  `json:"orderStatusString"`
	CreatedAt              time.Time               `json:"createdAt"`
	NewBeeMallOrderItemVOS []NewBeeMallOrderItemVO `json:"newBeeMallOrderItemVOS"`
}

type NewBeeMallOrderItemVO struct {
	GoodsId        int      `json:"goodsId"`
	GoodsName      string   `json:"goodsName"`
	GoodsCount     int      `json:"goodsCount"`
	GoodsCoverImg  string   `json:"goodsCoverImg"`
	SellingPrice   int      `json:"sellingPrice"`
	HashBlockChain []string `json:"hashBlockChain"`
}
