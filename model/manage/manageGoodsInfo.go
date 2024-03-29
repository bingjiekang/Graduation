// 自动生成模板MallGoodsInfo
package manage

import "time"

// MallGoodsInfo 结构体
type MallGoodsInfo struct {
	GoodsId            int       `json:"goodsId" form:"goodsId" gorm:"primarykey;AUTO_INCREMENT"`
	CommodityStocks    int       `json:"commodityStocks" form:"commodityStocks" gorm:"column:commodity_stocks;comment:商品库存位;"` // 商品库存位
	UUid               int64     `json:"uUid" form:"uUid" gorm:"column:u_uid;comment:唯一标识商家id"`
	GoodsName          string    `json:"goodsName" form:"goodsName" gorm:"column:goods_name;comment:商品名;type:varchar(200);"`
	GoodsIntro         string    `json:"goodsIntro" form:"goodsIntro" gorm:"column:goods_intro;comment:商品简介;type:varchar(200);"`
	GoodsCategoryId    int       `json:"goodsCategoryId" form:"goodsCategoryId" gorm:"column:goods_category_id;comment:关联分类id;type:bigint"`
	GoodsCoverImg      string    `json:"goodsCoverImg" form:"goodsCoverImg" gorm:"column:goods_cover_img;comment:商品主图;type:varchar(200);"`
	GoodsCarousel      string    `json:"goodsCarousel" form:"goodsCarousel" gorm:"column:goods_carousel;comment:商品轮播图;type:varchar(500);"`
	GoodsDetailContent string    `json:"goodsDetailContent" form:"goodsDetailContent" gorm:"column:goods_detail_content;comment:商品详情;type:text;"`
	OriginalPrice      int       `json:"originalPrice" form:"originalPrice" gorm:"column:original_price;comment:商品价格;type:int"`
	SellingPrice       int       `json:"sellingPrice" form:"sellingPrice" gorm:"column:selling_price;comment:商品实际售价;type:int"`
	StockNum           int       `json:"stockNum" form:"stockNum" gorm:"column:stock_num;comment:商品库存数量;type:int"`
	Tag                string    `json:"tag" form:"tag" gorm:"column:tag;comment:商品标签;type:varchar(20);"`
	GoodsSellStatus    int       `json:"goodsSellStatus" form:"goodsSellStatus" gorm:"column:goods_sell_status;comment:商品上架状态 1-下架 0-上架;type:tinyint"`
	PrevStock          int       `json:"prevStock" form:"prevStock" gorm:"column:prev_stock;comment:当前销售位;type:tinyint"`
	CreateUser         int64     `json:"createUser" form:"createUser" gorm:"column:create_user;comment:添加者主键id;type:int"`
	UpdateUser         int64     `json:"updateUser" form:"updateUser" gorm:"column:update_user;comment:修改者主键id;type:int"`
	CreatedAt          time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;type:datetime"`
	UpdatedAt          time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:更新时间;type:datetime"`
}

// TableName MallGoodsInfo 表名
func (MallGoodsInfo) TableName() string {
	return "mall_goods_info"
}
