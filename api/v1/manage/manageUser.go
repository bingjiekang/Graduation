package manage

import (
	"Graduation/global"
	"Graduation/model/common/request"
	"Graduation/model/common/response"
	req "Graduation/model/manage/request"
	"strconv"

	// "Graduation/model/manage/"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ManageAdminUserApi struct {
}

// 管理员用户登录(包括超级管理员)
// AdminLogin 管理员登陆
func (m *ManageAdminUserApi) ManageLogin(c *gin.Context) {
	var manageLoginParams req.ManageLoginParam
	_ = c.ShouldBindJSON(&manageLoginParams)
	if err, msg, token := manageUserService.ManageLogin(manageLoginParams); msg == "Ban" {
		response.FailWithMessage("抱歉,您已被禁用,请联系超级管理员解除!", c)
	} else if err != nil {
		response.FailWithMessage("登陆失败,请检查账号密码是否正确!", c)
	} else {
		response.OkWithData(token, c)
	}
}

// 管理员用户退出(包括超级管理员)
// AdminLogout 管理员登出
func (m *ManageAdminUserApi) ManageLogout(c *gin.Context) {
	token := c.GetHeader("token")
	if err := manageUserService.DeleteManageUserToken(token); err != nil {
		response.FailWithMessage("登出失败", c)
	} else {
		response.OkWithMessage("登出成功", c)
	}
}

// 管理员信息显示
// AdminUserProfile 用id查询AdminUser
func (m *ManageAdminUserApi) ManageUserInfo(c *gin.Context) {
	token := c.GetHeader("token")
	if err, mallAdminUser := manageUserService.GetManageUserInfo(token); err != nil {
		global.GVA_LOG.Error("未查询到管理员信息记录", zap.Error(err))
		response.FailWithMessage("未查询到管理员信息记录", c)
	} else {
		// 扰乱加密,防止泄露
		mallAdminUser.LoginPassword = "******"
		response.OkWithData(mallAdminUser, c)
	}
}

// 修改昵称
func (m *ManageAdminUserApi) UpdateManageUserNickName(c *gin.Context) {
	var reqs req.ManageUpdateNameParam
	_ = c.ShouldBindJSON(&reqs)
	token := c.GetHeader("token")
	if err := manageUserService.UpdateManageUserNickName(token, reqs); err != nil {
		global.GVA_LOG.Error("更新管理员用户昵称失败!", zap.Error(err))
		response.FailWithMessage("更新管理员用户昵称失败", c)
	} else {
		response.OkWithMessage("更新管理员用户昵称成功", c)
	}
}

// 修改密码
func (m *ManageAdminUserApi) UpdateManageUserPassword(c *gin.Context) {
	var reqs req.ManageUpdatePasswordParam
	_ = c.ShouldBindJSON(&reqs)
	userToken := c.GetHeader("token")
	if err := manageUserService.UpdateManagePassWord(userToken, reqs); err != nil {
		global.GVA_LOG.Error("更新密码失败!", zap.Error(err))
		response.FailWithMessage("更新密码失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("更新密码成功", c)
	}
}

// 用户商家列表显示
// UserList 商城注册商家用户列表
func (m *ManageAdminUserApi) UserList(c *gin.Context) {
	var pageInfo req.MallUserSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := manageUserService.GetManageUserInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取管理员用户失败!", zap.Error(err))
		response.FailWithMessage("获取管理员用户失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:       list,
			TotalCount: total,
			CurrPage:   pageInfo.PageNumber,
			PageSize:   pageInfo.PageSize,
		}, "获取管理员用户成功", c)
	}
}

// LockUser 用户禁用[0]与解除禁用[1](0-未锁定 1-已锁定)
func (m *ManageAdminUserApi) LockUser(c *gin.Context) {
	lockStatus, _ := strconv.Atoi(c.Param("lockStatus"))
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := manageUserService.LockUser(IDS, lockStatus); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}
