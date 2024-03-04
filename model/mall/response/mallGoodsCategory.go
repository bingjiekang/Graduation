package response

// 首页分类数据VO(第三级)
type ThirdLevelCategoryVO struct {
	CategoryId    int    `json:"categoryId"`
	CategoryLevel int    `json:"categoryLevel"`
	CategoryName  string `json:"categoryName" `
}

// 第二级分类数据
type SecondLevelCategoryVO struct {
	CategoryId            int                    `json:"categoryId"`
	ParentId              int                    `json:"parentId"`
	CategoryLevel         int                    `json:"categoryLevel"`
	CategoryName          string                 `json:"categoryName" `
	ThirdLevelCategoryVOS []ThirdLevelCategoryVO `json:"thirdLevelCategoryVOS"`
}

// 分类页左侧第一级显示数据
type MallIndexCategoryVO struct {
	CategoryId int `json:"categoryId"`
	//ParentId               int                      `json:"parentId"`
	CategoryLevel          int                     `json:"categoryLevel"`
	CategoryName           string                  `json:"categoryName" `
	SecondLevelCategoryVOS []SecondLevelCategoryVO `json:"secondLevelCategoryVOS"`
}
