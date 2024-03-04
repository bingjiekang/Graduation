package mall

import "time"

// 轮播图配置
type MallCarousel struct {
	CarouselId   int       `json:"carouselId" form:"carouselId" gorm:"primarykey;AUTO_INCREMENT"`
	CarouselUrl  string    `json:"carouselUrl" form:"carouselUrl" gorm:"column:carousel_url;comment:轮播图;type:varchar(100);"`
	RedirectUrl  string    `json:"redirectUrl" form:"redirectUrl" gorm:"column:redirect_url;comment:点击后的跳转地址(默认不跳转);type:varchar(100);"`
	CarouselRank int       `json:"carouselRank" form:"carouselRank" gorm:"column:carousel_rank;comment:排序值(字段越大越靠前);type:int"`
	IsDeleted    int       `json:"isDeleted" form:"isDeleted" gorm:"column:is_deleted;comment:删除标识字段(0-未删除 1-已删除);type:tinyint"`
	CreateUser   int       `json:"createUser" form:"createUser" gorm:"column:create_user;comment:创建者id;type:int"`
	UpdateUser   int       `json:"updateUser" form:"updateUser" gorm:"column:update_user;comment:修改者id;type:int"`
	CreatedAt    time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;type:datetime"`
	UpdatedAt    time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:更新时间;type:datetime"`
}

// TableName MallCarousel 表名
func (MallCarousel) TableName() string {
	return "mall_carousel"
}
