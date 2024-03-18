package manage

import "time"

// 用户订单表(记得和用户订单表一起更新)
type MallAdminOrder struct {
	OrderId     int       `json:"orderId" form:"orderId" gorm:"primarykey;AUTO_INCREMENT"` // 自增id
	OrderNo     string    `json:"orderNo" form:"orderNo" gorm:"column:order_no;comment:订单号;type:varchar(20);"`
	Muid        int64     `json:"mUid" form:"mUid" gorm:"column:m_uid;comment:卖家id"`
	Buid        int64     `json:"bUid" form:"bUid" gorm:"column:b_uid; comment:买家uid"`
	TotalPrice  int       `json:"totalPrice" form:"totalPrice" gorm:"column:total_price;comment:订单总价;type:int"`
	PayStatus   int       `json:"payStatus" form:"payStatus" gorm:"column:pay_status;comment:支付状态:0.未支付,1.支付成功,-1:支付失败;type:tinyint"`
	PayType     int       `json:"payType" form:"payType" gorm:"column:pay_type;comment:0.无 1.支付宝支付 2.微信支付;type:tinyint"`
	PayTime     time.Time `json:"payTime" form:"payTime" gorm:"column:pay_time;comment:支付时间;type:datetime"`
	OrderStatus int       `json:"orderStatus" form:"orderStatus" gorm:"column:order_status;comment:订单状态:0.待支付 1.已支付 2.配货完成 3:出库成功 4.交易成功 -1.手动关闭 -2.超时关闭 -3.商家关闭;type:tinyint"`
	ExtraInfo   string    `json:"extraInfo" form:"extraInfo" gorm:"column:extra_info;comment:订单body;type:varchar(100);"`
	IsDeleted   int       `json:"isDeleted" form:"isDeleted" gorm:"column:is_deleted;comment:删除标识字段(0-未删除 1-已删除);type:tinyint"`
	CreatedAt   time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;type:datetime"`
	UpdatedAt   time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:更新时间;type:datetime"`
}

// TableName MallAdminOrder 表名
func (MallAdminOrder) TableName() string {
	return "mall_admin_order"
}
