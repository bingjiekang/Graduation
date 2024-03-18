package manage

import "time"

// 区块链商品交易信息
type MallBlockTrading struct {
	Id              int64     `json:"id" form:"id" gorm:"primarykey;AUTO_INCREMENT;"`                                       // 主键自增id
	OrderNo         string    `json:"orderNo" form:"orderNo" gorm:"column:order_no;comment:订单号;type:varchar(20);"`          // 订单号
	Commodity       int       `json:"commodity" form:"commodity" gorm:"column:commodity;comment:商品序号;"`                     // 商品序号
	CommodityStocks int       `json:"commodityStocks" form:"commodityStocks" gorm:"column:commodity_stocks;comment:商品库存位;"` // 商品库存位
	SellerUid       int64     `json:"sellerUid" form:"sellerUid" gorm:"column:seller_uid;comment:售卖者商家uid;"`                // 商家uid
	BuyerUid        int64     `json:"buyerUid" form:"buyerUid" gorm:"column:buyer_uid;comment:买家用户uid;"`                    // 用户购买者uid
	InitBlockHash   string    `json:"initBlockHash" form:"initBlockHash" gorm:"column:init_block_hash;comment:初始商品区块哈希;"`   // 初始商品区块哈希
	CurrBlockHash   string    `json:"currBlockHash" form:"currBlockHash" gorm:"column:curr_block_hash;comment:当前商品区块哈希;"`   // 当前商品区块哈希
	CreatedAt       time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;type:datetime"`       // 创建时间
	UpdatedAt       time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:更新时间;type:datetime"`       // 更新时间
}

// 商品交易信息数据库命名
func (MallBlockTrading) TableName() string {
	return "block_trading"
}
