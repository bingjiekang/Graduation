package manage

import (
	v1 "Graduation/api/v1"
	"Graduation/middleware"

	"github.com/gin-gonic/gin"
)

type ManageAdminUserRouter struct {
}

func (r *ManageAdminUserRouter) ApiManageAdminUserRouter(Router *gin.RouterGroup) {
	mallAdminUserWithoutRouter := Router.Group("v1")                                                              // 不验证权限
	mallAdminUserRouter := Router.Group("v1").Use(middleware.AdminJWTAuth())                                      // 验证管理员权限
	mallSuperAdminUserRouter := Router.Group("v1").Use(middleware.AdminJWTAuth(), middleware.SuperAdminJWTAuth()) // 验证超级管理员权限

	var mallAdminUserApi = v1.ApiGroupApp.ManageApiGroup.ManageAdminUserApi
	// 普通无权限
	{
		mallAdminUserWithoutRouter.POST("adminUser/login", mallAdminUserApi.ManageLogin) //管理员登陆
		mallAdminUserWithoutRouter.POST("upload/file", mallAdminUserApi.UploadFile)      // 上传图片
		mallAdminUserWithoutRouter.POST("upload/files", mallAdminUserApi.UploadFiles)    // 上传多张图片
	}
	// 管理员权限
	{

		mallAdminUserRouter.DELETE("logout", mallAdminUserApi.ManageLogout)                      // 登出
		mallAdminUserRouter.GET("adminUser/profile", mallAdminUserApi.ManageUserInfo)            // 根据ID获取 管理员详情(用来显示信息)
		mallAdminUserRouter.PUT("adminUser/name", mallAdminUserApi.UpdateManageUserNickName)     // 更新管理员用户昵称
		mallAdminUserRouter.PUT("adminUser/password", mallAdminUserApi.UpdateManageUserPassword) // 更新管理员用户密码

		// mallAdminUserRouter.POST("createMallAdminUser", mallAdminUserApi.CreateAdminUser) // 新建MallAdminUser
		// mallAdminUserRouter.POST("upload/file", mallAdminUserApi.UploadFile) //上传图片

	}
	// 超级管理员权限
	{
		mallSuperAdminUserRouter.GET("users", mallAdminUserApi.UserList)             // 获取管理员用户信息列表
		mallSuperAdminUserRouter.PUT("users/:lockStatus", mallAdminUserApi.LockUser) // 锁定解锁管理员
	}
}
