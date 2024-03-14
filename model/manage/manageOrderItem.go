package manage

import "time"

type MallOrderItem struct {
	OrderItemId     int       `json:"orderItemId" gorm:"primarykey;AUTO_INCREMENT"`
	OrderId         int       `json:"orderId" form:"orderId" gorm:"column:order_id;;type:bigint"`
	GoodsId         int       `json:"goodsId" form:"goodsId" gorm:"column:goods_id;;type:bigint"`
	CommodityStocks int64     `json:"commodityStocks" form:"commodityStocks" gorm:"column:commodity_stocks;comment:商品库存位;"` // 商品库存位
	GoodsName       string    `json:"goodsName" form:"goodsName" gorm:"column:goods_name;comment:商品名;type:varchar(200);"`
	GoodsCoverImg   string    `json:"goodsCoverImg" form:"goodsCoverImg" gorm:"column:goods_cover_img;comment:商品主图;type:varchar(200);"`
	SellingPrice    int       `json:"sellingPrice" form:"sellingPrice" gorm:"column:selling_price;comment:商品实际售价;type:int"`
	GoodsCount      int       `json:"goodsCount" form:"goodsCount" gorm:"column:goods_count;;type:bigint"`
	CreatedAt       time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;type:datetime"`
}

func (MallOrderItem) TableName() string {
	return "mall_order_item"
}
