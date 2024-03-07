package service

import (
	"Graduation/service/mall"
	"Graduation/service/manage"
)

type ServiceGroup struct {
	MallServiceGroup   mall.MallServiceGroup
	ManageServiceGroup manage.ManageServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
