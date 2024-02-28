package service

import "Graduation/service/mall"

type ServiceGroup struct {
	MallServiceGroup mall.MallServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
