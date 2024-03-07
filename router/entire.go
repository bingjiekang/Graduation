package router

import (
	"Graduation/router/mall"
	"Graduation/router/manage"
)

// api接口组
type RouterGroup struct {
	// Manage manage.ManageRouterGroup
	Mall   mall.MallRouterGroup
	Manage manage.ManageRouterGroup
}

var RouterGroupApp = new(RouterGroup)
