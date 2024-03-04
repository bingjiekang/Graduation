package mall

import "time"

// 购物车列表信息 结构体
type MallShopCartItem struct {
	CartItemId int       `json:"cartItemId" form:"cartItemId" gorm:"primarykey;AUTO_INCREMENT"`
	UUid       int64     `json:"uUid" form:"uUid" gorm:"column:u_uid;comment:唯一标识id"`
	GoodsId    int       `json:"goodsId" form:"goodsId" gorm:"column:goods_id;comment:关联商品id;type:bigint"`
	GoodsCount int       `json:"goodsCount" form:"goodsCount" gorm:"column:goods_count;comment:数量(最大为5);type:int"`
	IsDeleted  int       `json:"isDeleted" form:"isDeleted" gorm:"column:is_deleted;comment:删除标识字段(0-未删除 1-已删除);type:tinyint"`
	CreatedAt  time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;type:datetime"`
	UpdatedAt  time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:更新时间;type:datetime"`
}

// TableName MallShoppingCartItem 表名
func (MallShopCartItem) TableName() string {
	return "mall_shop_cart_item"
}
