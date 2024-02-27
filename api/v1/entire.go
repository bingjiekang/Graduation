package v1

import "Graduation/api/v1/mall"

type ApiGroup struct {
	MallApiGroup mall.MallGroup
}

// 创建一个新的用户api组
var ApiGroupApp = new(ApiGroup)
