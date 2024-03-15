package mall

import (
	"Graduation/global"
	"Graduation/model/common/response"
	"Graduation/model/mall/request"
	"Graduation/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MallUserApi struct {
}

// 处理用户注册的路由接口转接
func (m *MallUserApi) UserRegister(c *gin.Context) {
	// 绑定对应用户注册信息结构体
	var req request.RegisterUserParam
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("注册信息无法绑定对应结构")
	}
	// 检查用户传入的信息是否合理
	if !(utils.ValidatePhoneNumber(req.LoginName) && utils.ValidatePassword(req.Password)) { // 有一个不合理
		response.FailWithMessage("请确保用户名为手机号,密码为8位以上数字,密码,特殊符合的组合!", c)
		return
	}
	// 对用户进行检查和注册处理
	if err := mallUserService.RegisterUser(req); err != nil {
		global.GVA_LOG.Error("创建失败", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功!", c)
}

// 处理用户登录的路由接口转接
func (m *MallUserApi) UserLogin(c *gin.Context) {
	// 绑定对应登录信息结构体
	var req request.UserLoginParam
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("登陆信息无法绑定对应结构")
	}
	// 校验登陆信息是否正确
	if err, _, token := mallUserService.LoginUser(req); err != nil {
		response.FailWithMessage("登陆失败!请检查账号和密码是否错误", c)
	} else {
		response.OkWithData(token, c)
	}
}

// 处理用户退出的路由接口转接
func (m *MallUserApi) UserLogout(c *gin.Context) {
	token := c.GetHeader("token")
	global.GVA_LOG.Info("登陆token:" + token)
	// 检查并删除用户 token
	if err := mallUserService.DeleteMallUserToken(token); err != nil {
		response.FailWithMessage("登出失败", c)
	} else {
		response.OkWithMessage("登出成功", c)
	}
}

// 处理用户信息的路由接口转接
func (m *MallUserApi) UserInfo(c *gin.Context) {
	token := c.GetHeader("token")
	// 获取用户信息并返回
	if err, userDetail := mallUserService.GetUserInfo(token); err != nil {
		global.GVA_LOG.Error("未查询到用户记录", zap.Error(err))
		response.FailWithMessage("未查询到用户记录", c)
	} else {
		response.OkWithData(userDetail, c)
	}
}

// 用户修改信息接口路由
func (m *MallUserApi) UpdateUserInfo(c *gin.Context) {
	token := c.GetHeader("token")
	var req request.UpdateUserInfoParam
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("用户更新的信息无法绑定对应结构")
	}
	// 获取用户信息
	if err := mallUserService.UpdateUserInfo(token, req); err != nil {
		global.GVA_LOG.Error("更新用户信息失败", zap.Error(err))
		response.FailWithMessage("更新用户信息失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
