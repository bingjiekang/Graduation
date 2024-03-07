package v1

import (
	"Graduation/api/v1/mall"
	"Graduation/api/v1/manage"
)

type ApiGroup struct {
	MallApiGroup   mall.MallGroup
	ManageApiGroup manage.ManageGroup
}

// 创建一个新的用户api组
var ApiGroupApp = new(ApiGroup)
