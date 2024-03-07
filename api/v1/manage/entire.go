package manage

import "Graduation/service"

type ManageGroup struct {
	ManageAdminUserApi
}

var (
	manageUserService = service.ServiceGroupApp.ManageServiceGroup.ManageUserService // 管理员以及超级管理员操作
)
