package mall

// 用户订单地址
type MallOrderAddress struct {
	OrderId       int    `json:"orderId"`
	UserName      string `json:"userName"`
	UserPhone     string `json:"userPhone"`
	ProvinceName  string `json:"provinceName"`
	CityName      string `json:"cityName"`
	RegionName    string `json:"regionName"`
	DetailAddress string `json:"detailAddress"`
}
