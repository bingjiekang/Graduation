package mall

import "time"

// 种类分类表
type MallGoodsCategory struct {
	CategoryId    int       `json:"categoryId" gorm:"primarykey;AUTO_INCREMENT"`
	CategoryLevel int       `json:"categoryLevel" gorm:"comment:分类等级"`
	ParentId      int       `json:"parentId" gorm:"comment:父类id"`
	CategoryName  string    `json:"categoryName" gorm:"comment:分类名称"`
	CategoryRank  int       `json:"categoryRank" gorm:"comment:排序比重"`
	IsDeleted     int       `json:"isDeleted" gorm:"comment:是否删除"`
	CreatedAt     time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;type:datetime"`
	UpdatedAt     time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:更新时间;type:datetime"`
}

func (MallGoodsCategory) TableName() string {
	return "mall_goods_category"
}
