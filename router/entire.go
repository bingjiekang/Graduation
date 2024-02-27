package router

import "Graduation/router/mall"

// api接口组
type RouterGroup struct {
	// Manage manage.ManageRouterGroup
	Mall mall.MallRouterGroup
}

var RouterGroupApp = new(RouterGroup)
