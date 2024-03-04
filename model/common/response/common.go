package response

// 返回对应页数的数据
type PageResult struct {
	List       interface{} `json:"list"`
	TotalCount int64       `json:"totalCount"`
	TotalPage  int         `json:"totalPage"`
	CurrPage   int         `json:"currPage"`
	PageSize   int         `json:"pageSize"`
}
