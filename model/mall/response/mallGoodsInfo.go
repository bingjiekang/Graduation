package response

// 商品详情信息界面
type GoodsInfoDetailResponse struct {
	GoodsId            int      `json:"goodsId"`
	GoodsName          string   `json:"goodsName"`
	GoodsIntro         string   `json:"goodsIntro"`
	GoodsCoverImg      string   `json:"goodsCoverImg"`
	SellingPrice       int      `json:"sellingPrice"`
	GoodsDetailContent string   `json:"goodsDetailContent"  `
	OriginalPrice      int      `json:"originalPrice" `
	Tag                string   `json:"tag" form:"tag" `
	GoodsCarouselList  []string `json:"goodsCarouselList" `
}

// 商品搜索信息界面
type GoodsSearchResponse struct {
	GoodsId       int    `json:"goodsId"`
	GoodsName     string `json:"goodsName"`
	GoodsIntro    string `json:"goodsIntro"`
	GoodsCoverImg string `json:"goodsCoverImg"`
	SellingPrice  int    `json:"sellingPrice"`
}
