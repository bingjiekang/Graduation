package mall

import "time"

// 用户地址
type MallUserAddress struct {
	AddressId     int       `json:"addressId" form:"addressId" gorm:"primarykey;AUTO_INCREMENT"`
	Uuid          int64     `json:"uUid" form:"uUid" gorm:"column:u_uid;comment:uid"`
	UserName      string    `json:"userName" form:"userName" gorm:"column:user_name;comment:收货人姓名;type:varchar(30);"`
	UserPhone     string    `json:"userPhone" form:"userPhone" gorm:"column:user_phone;comment:收货人手机号;type:varchar(11);"`
	DefaultFlag   int       `json:"defaultFlag" form:"defaultFlag" gorm:"column:default_flag;comment:是否为默认 0-非默认 1-是默认;type:tinyint"`
	ProvinceName  string    `json:"provinceName" form:"provinceName" gorm:"column:province_name;comment:省;type:varchar(32);"`
	CityName      string    `json:"cityName" form:"cityName" gorm:"column:city_name;comment:城;type:varchar(32);"`
	RegionName    string    `json:"regionName" form:"regionName" gorm:"column:region_name;comment:区;type:varchar(32);"`
	DetailAddress string    `json:"detailAddress" form:"detailAddress" gorm:"column:detail_address;comment:收件详细地址(街道/楼宇/单元);type:varchar(64);"`
	IsDeleted     int       `json:"isDeleted" form:"isDeleted" gorm:"column:is_deleted;comment:删除标识字段(0-未删除 1-已删除);type:tinyint"`
	CreatedAt     time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;type:datetime"`
	UpdatedAt     time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:更新时间;type:datetime"`
}

// TableName MallUserAddress 表名
func (MallUserAddress) TableName() string {
	return "mall_user_address"
}
