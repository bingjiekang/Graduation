package request

// 添加用户地址请求结构体
type AddAddressParam struct {
	UserName      string `json:"userName"`
	UserPhone     string `json:"userPhone"`
	DefaultFlag   byte   `json:"defaultFlag"` // 0-不是 1-是
	ProvinceName  string `json:"provinceName"`
	CityName      string `json:"cityName"`
	RegionName    string `json:"regionName"`
	DetailAddress string `json:"detailAddress"`
}

// 更新用户地址结构体
type UpdateAddressParam struct {
	AddressId     string `json:"addressId"`
	Uuid          int64  `json:"uUid"`
	UserName      string `json:"userName"`
	UserPhone     string `json:"userPhone"`
	DefaultFlag   byte   `json:"defaultFlag"` // 0-不是 1-是
	ProvinceName  string `json:"provinceName"`
	CityName      string `json:"cityName"`
	RegionName    string `json:"regionName"`
	DetailAddress string `json:"detailAddress"`
}
