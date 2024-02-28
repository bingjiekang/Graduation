package mall

import (
	"Graduation/global"
	"Graduation/model/common/response"
	"Graduation/model/mall/request"
	"Graduation/utils"
	"fmt"

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
		response.FailWithMessage("登陆失败!", c)
	} else {
		response.OkWithData(token, c)
	}
	// // fmt.Println("接收登陆")
	// response.OkWithMessage("登陆成功!", c)

}

// 处理用户退出的路由接口转接
func (m *MallUserApi) UserLogout(c *gin.Context) {
	token := c.GetHeader("token")
	fmt.Println("登陆token", token)
	if err := mallUserService.DeleteMallUserToken(token); err != nil {
		response.FailWithMessage("登出失败", c)
	} else {
		response.OkWithMessage("登出成功", c)
	}
	response.OkWithMessage("登出成功", c)

}
