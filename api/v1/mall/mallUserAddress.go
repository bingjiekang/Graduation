package mall

import (
	"Graduation/global"
	"Graduation/model/common/response"
	"Graduation/model/mall/request"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MallUserAddressApi struct {
}

// 增加用户地址信息
func (m *MallUserAddressApi) AddUserAddress(c *gin.Context) {
	var req request.AddAddressParam
	_ = c.ShouldBindJSON(&req)
	token := c.GetHeader("token")
	// 保存用户地址信息
	err := mallUserAddressService.AddUserAddress(token, req)
	if err != nil {
		global.GVA_LOG.Error("创建失败", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// 查询用户地址列表信息
func (m *MallUserAddressApi) GetAddressList(c *gin.Context) {
	token := c.GetHeader("token")
	if err, userAddressList := mallUserAddressService.GetUserAddressList(token); err != nil {
		global.GVA_LOG.Error("获取地址列表信息失败", zap.Error(err))
		response.FailWithMessage("获取地址列表信息失败:"+err.Error(), c)
	} else if len(userAddressList) == 0 {
		global.GVA_LOG.Info("获取地址列表信息为空")
		response.OkWithData(nil, c)
	} else {
		response.OkWithData(userAddressList, c)
	}
}

// 查看用户指定标识的地址信息
func (m *MallUserAddressApi) GetUserAddress(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("addressId"), 10, 64)
	token := c.GetHeader("token")
	if err, userAddress := mallUserAddressService.GetUserAddress(token, id); err != nil {
		global.GVA_LOG.Error("获取指定地址信息失败", zap.Error(err))
		response.FailWithMessage("获取指定地址信息失败:"+err.Error(), c)
	} else {
		response.OkWithData(userAddress, c)
	}
}

// 修改用户地址信息
func (m *MallUserAddressApi) UpdateUserAddress(c *gin.Context) {
	// 接受修改的用户信息并绑定对应结构体
	var req request.UpdateAddressParam
	_ = c.ShouldBindJSON(&req)
	token := c.GetHeader("token")
	// 修改用户地址信息
	err := mallUserAddressService.UpdateUserAddress(token, req)
	if err != nil {
		global.GVA_LOG.Error("用户地址信息修改失败", zap.Error(err))
		response.FailWithMessage("用户地址信息修改失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("用户地址信息修改成功", c)
	}
}

// 删除指定用户标识的地址信息
func (m *MallUserAddressApi) DeleteUserAddress(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("addressId"), 10, 64)
	token := c.GetHeader("token")
	if err := mallUserAddressService.DeleteUserAddress(token, id); err != nil {
		global.GVA_LOG.Error("删除用户指定地址信息失败", zap.Error(err))
		response.FailWithMessage("删除用户指定地址信息失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("删除用户地址成功", c)
	}
}
