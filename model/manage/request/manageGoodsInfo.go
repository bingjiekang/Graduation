package request

import (
	"Graduation/model/common/request"
	"Graduation/model/manage"
	"time"
)

type MallGoodsInfoSearch struct {
	manage.MallGoodsInfo
	request.PageInfo
}

// 添加商品
type GoodsInfoAddParam struct {
	GoodsName          string      `json:"goodsName"`
	GoodsIntro         string      `json:"goodsIntro"`
	GoodsCategoryId    int         `json:"goodsCategoryId"`
	GoodsCoverImg      string      `json:"goodsCoverImg"`
	GoodsCarousel      string      `json:"goodsCarousel"`
	GoodsDetailContent string      `json:"goodsDetailContent"`
	OriginalPrice      interface{} `json:"originalPrice"`
	SellingPrice       interface{} `json:"sellingPrice"`
	StockNum           interface{} `json:"stockNum"`
	Tag                string      `json:"tag"`
	GoodsSellStatus    string      `json:"goodsSellStatus"`
}

// GoodsInfoUpdateParam 更新商品信息的入参
type GoodsInfoUpdateParam struct {
	GoodsId            string      `json:"goodsId"`
	GoodsName          string      `json:"goodsName"`
	GoodsIntro         string      `json:"goodsIntro"`
	GoodsCategoryId    int         `json:"goodsCategoryId"`
	GoodsCoverImg      string      `json:"goodsCoverImg"`
	GoodsCarousel      string      `json:"goodsCarousel"`
	GoodsDetailContent string      `json:"goodsDetailContent"`
	OriginalPrice      interface{} `json:"originalPrice"`
	SellingPrice       interface{} `json:"sellingPrice"`
	StockNum           interface{} `json:"stockNum"`
	Tag                string      `json:"tag"`
	GoodsSellStatus    int         `json:"goodsSellStatus"`
	CreatedAt          time.Time   `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;type:datetime"`
	UpdatedAt          time.Time   `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:更新时间;type:datetime"`
}

type StockNumDTO struct {
	GoodsId    int `json:"goodsId"`
	GoodsCount int `json:"goodsCount"`
}
