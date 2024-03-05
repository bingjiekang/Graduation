package mall

import (
	v1 "Graduation/api/v1"
	"Graduation/middleware"

	"github.com/gin-gonic/gin"
)

// 用户地址路由连接表
type MallUserAddressRouter struct {
}

func (m *MallUserAddressRouter) ApiMallUserAddressRouter(Router *gin.RouterGroup) {
	mallUserAddressRouter := Router.Group("v1").Use(middleware.UserJWTAuth())
	var userAddressApi = v1.ApiGroupApp.MallApiGroup.MallUserAddressApi
	{
		mallUserAddressRouter.POST("/address", userAddressApi.AddUserAddress)                 // 增加地址
		mallUserAddressRouter.GET("/address", userAddressApi.GetAddressList)                  // 查看用户全部地址列表信息
		mallUserAddressRouter.GET("/address/:addressId", userAddressApi.GetUserAddress)       // 获取指定地址详情
		mallUserAddressRouter.PUT("/address", userAddressApi.UpdateUserAddress)               // 修改用户指定地址信息
		mallUserAddressRouter.DELETE("/address/:addressId", userAddressApi.DeleteUserAddress) //删除地址
		mallUserAddressRouter.GET("/address/default", userAddressApi.GetUserDefaultAddress)   //获取默认地址

	}

}
