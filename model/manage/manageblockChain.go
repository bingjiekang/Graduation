package manage

import "time"

// 区块链商品哈希
type MallBlockChain struct {
	Id              int64     `json:"id" form:"id" gorm:"primarykey;AUTO_INCREMENT;"`                                       // 主键自增id
	Commodity       int64     `json:"commodity" form:"commodity" gorm:"column:commodity;comment:商品序号;"`                     // 商品序号
	CommodityStocks int64     `json:"commodityStocks" form:"commodityStocks" gorm:"column:commodity_stocks;comment:商品库存位;"` // 商品库存位
	IsSale          bool      `json:"isSale" form:"isSale" gorm:"column:is_sale;comment:是否出售默认为false;"`                     // 是否出售
	PduIdent        string    `json:"pduIdent" form:"pduIdent" gorm:"column:pdu_ident;comment:商品唯一标识;type:varchar(200);"`   // 商品的唯一标识信息(即初始创世区块信息)
	Number          int64     `json:"number" form:"number" gorm:"column:number;comment:交易次数;"`                              // 商品交易次数
	InitBlockHash   string    `json:"initBlockHash" form:"initBlockHash" gorm:"column:init_block_hash;comment:初始商品区块哈希;"`   // 初始商品区块哈希
	CurrBlockHash   string    `json:"currBlockHash" form:"currBlockHash" gorm:"column:curr_block_hash;comment:当前商品区块哈希;"`   // 当前商品区块哈希
	CreatedAt       time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;type:datetime"`       // 创建时间
	UpdatedAt       time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:更新时间;type:datetime"`       // 更新时间
}

func (MallBlockChain) TableName() string {
	return "block_chain"
}
