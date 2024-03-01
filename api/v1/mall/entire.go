package mall

import "Graduation/service"

type MallGroup struct {
	MallUserApi
	MallUserAddressApi
}

var mallUserService = service.ServiceGroupApp.MallServiceGroup.MallUserService
var mallUserAddressService = service.ServiceGroupApp.MallServiceGroup.MallUserAddressService
