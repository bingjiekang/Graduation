package mall

import "Graduation/service"

type MallGroup struct {
	MallUserApi
}

var mallUserService = service.ServiceGroupApp.MallServiceGroup.MallUserService
