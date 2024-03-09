package request

import (
	"Graduation/model/common/request"
	"time"
)

// 添加商品信息请求
type MallGoodsCategoryReq struct {
	CategoryId    int       `json:"categoryId"`
	CategoryLevel string    `json:"categoryLevel" `
	ParentId      string    `json:"parentId"`
	CategoryName  string    `json:"categoryName" `
	CategoryRank  string    `json:"categoryRank" `
	IsDeleted     int       `json:"isDeleted" `
	CreatedAt     time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;type:datetime"`
	UpdatedAt     time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:更新时间;type:datetime"`
}

// 查询商品信息分类
type SearchCategoryParams struct {
	CategoryLevel int `json:"categoryLevel" form:"categoryLevel"`
	ParentId      int `json:"parentId" form:"parentId"`
	request.PageInfo
}
