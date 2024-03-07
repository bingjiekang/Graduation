package manage

import (
	"Graduation/global"
	"Graduation/model/common/request"
	"Graduation/model/mall"
	mag "Graduation/model/manage"
	req "Graduation/model/manage/request"
	mallservice "Graduation/service/mall"
	"Graduation/utils"
	"errors"
	"fmt"
	"time"
)

type ManageUserService struct {
}

// 管理员用户登录以及超级管理员登录
func (m *ManageUserService) ManageLogin(req req.ManageLoginParam) (err error, msg string, token string) {
	var mallAdminUser mag.MallAdminUser
	// 管理员就是用户(先查用户表,用户表不存在,则无法登陆,用户表存在则加入管理员账户)
	var user mall.MallUser
	err = global.GVA_DB.Where("login_name=? AND password_md5=?", req.UserName, req.PasswordMd5).First(&user).Error
	if err != nil || user == (mall.MallUser{}) {
		return err, "不存在用户,请注册成为商城用户后才有权限登录", token
	}
	// 查询到用户已存在,
	// 判断是否已经加入到管理员表
	err = global.GVA_DB.Where("login_user_name=? AND login_password=?", req.UserName, req.PasswordMd5).First(&mallAdminUser).Error
	if mallAdminUser == (mag.MallAdminUser{}) {
		// 没加入 则将用户信息加入到管理员信息表super_admin_user
		{
			mallAdminUser.UUid = user.UUid
			mallAdminUser.LoginUserName = user.LoginName
			mallAdminUser.LoginPassword = user.PasswordMd5
			mallAdminUser.NickName = user.NickName
			mallAdminUser.IsSuperAdmin = 0 // 普通管理员

		}
	}
	// 更新管理员表和用户表锁定状态相同
	mallAdminUser.Locked = user.LockedFlag
	// 如果用户被禁,则无法登陆后台管理员系统
	if mallAdminUser.Locked == 1 {
		// 但仍然需要更新用户信息数据
		err = global.GVA_DB.Save(&mallAdminUser).Error
		return err, "Ban", token
	}
	// 创建token
	token, _ = utils.CreateToken(user.UUid)
	strUuid := fmt.Sprintf("%d", user.UUid)
	err = global.GVA_REDIS.Set(global.GVA_CTX, strUuid, token, 3600*time.Second).Err()
	if err != nil {
		global.GVA_LOG.Error("redis存储token失败")
		return err, msg, token
	} else {
		global.GVA_LOG.Info("redis存储token成功")
	}
	err = global.GVA_DB.Save(&mallAdminUser).Error
	return

}

// 管理员以及超级管理员登出,删除管理员登陆token
func (m *ManageUserService) DeleteManageUserToken(token string) (err error) {
	uuid, err, ok := utils.UndoToken(token)
	if err != nil && ok == 0 { // 解码token出现错误
		return err
	}
	uid := fmt.Sprintf("%d", uuid)
	if ok == 1 { // 超时 token已经失效
		global.GVA_REDIS.Del(global.GVA_CTX, uid)
		return nil
	}
	_, err = global.GVA_REDIS.Get(global.GVA_CTX, uid).Result()
	if err != nil {
		global.GVA_LOG.Error("无法从redis得到对应uid的token信息,可能不存在")
		return err
	}
	// 存在则删除
	global.GVA_REDIS.Del(global.GVA_CTX, uid)
	return nil
}

// 检查管理员token是否存在
func (m *ManageUserService) ExistManageToken(token string) (err error, tm int64) {
	uuid, err, ok := utils.UndoToken(token)
	if err != nil && ok == 0 { // 解码token出现错误
		return err, 0
	}
	uid := fmt.Sprintf("%d", uuid)
	_, err = global.GVA_REDIS.Get(global.GVA_CTX, uid).Result()
	if err != nil {
		global.GVA_LOG.Error("token不存在")
		return err, 0
	}
	if ok == 1 {
		global.GVA_LOG.Info("token已超时")
		return err, 1
	}
	return nil, 2
}

// 获取登录管理员和超级管理员的信息
func (m *ManageUserService) GetManageUserInfo(token string) (err error, mallAdminUser mag.MallAdminUser) {
	// 判断用户是否存在
	if !mallservice.IsUserExist(token) {
		return errors.New("不存在的用户"), mallAdminUser
	}
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid = ?", uuid).First(&mallAdminUser).Error
	return err, mallAdminUser
}

// 判断是否是超级管理员
func (m *ManageUserService) IsSuperManageAdmin(token string) (err error, ok bool) {
	// 判断用户是否存在
	if !mallservice.IsUserExist(token) {
		return errors.New("不存在的用户"), false
	}
	var mallAdminUser mag.MallAdminUser
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid = ? And is_super_admin = 1", uuid).First(&mallAdminUser).Error
	if err != nil {
		return err, false
	}
	// 不为空结构体,证明查询到内容
	if mallAdminUser != (mag.MallAdminUser{}) {
		return nil, true
	}
	return nil, false
}

// 更新管理员用户昵称
func (m *ManageUserService) UpdateManageUserNickName(token string, reqs req.ManageUpdateNameParam) (err error) {
	// 判断用户是否存在
	if !mallservice.IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	uuid, _, _ := utils.UndoToken(token)
	// 更新管理员表的昵称
	if err = global.GVA_DB.Where("u_uid = ?", uuid).Updates(&mag.MallAdminUser{
		NickName: reqs.NickName,
	}).Error; err != nil {
		return err
	}

	// 更新用户表里的昵称
	if err = global.GVA_DB.Where("u_uid = ?", uuid).Updates(&mall.MallUser{
		NickName: reqs.NickName,
	}).Error; err != nil {
		return err
	}
	return

}

// 更新管理员用户密码
func (m *ManageUserService) UpdateManagePassWord(token string, reqs req.ManageUpdatePasswordParam) (err error) {
	// 判断用户是否存在
	if !mallservice.IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	var adminUser mag.MallAdminUser
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid =?", uuid).First(&adminUser).Error
	if err != nil {
		return errors.New("不存在的管理员用户")
	}
	if adminUser.LoginPassword != reqs.OriginalPassword {
		return errors.New("原密码不正确")
	}
	if reqs.NewPassword == "" {
		return errors.New("密码不能为空")
	}
	adminUser.LoginPassword = reqs.NewPassword
	// 更新管理员表
	if err = global.GVA_DB.Where("u_uid=?", uuid).Updates(&adminUser).Error; err != nil {
		return
	}
	// 更新用户表
	var userInfo mall.MallUser
	if err = global.GVA_DB.Where("u_uid =?", uuid).First(&userInfo).Error; err != nil {
		return errors.New("管理员用户密码更新失败")
	}
	userInfo.PasswordMd5 = reqs.NewPassword
	err = global.GVA_DB.Save(&userInfo).Error
	return

}

// 查看用户管理员的信息列表
// GetManageUserInfoList 分页获取商城注册用户即管理员列表
func (m *ManageUserService) GetManageUserInfoList(info req.MallUserSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageNumber - 1)
	// 创建db
	db := global.GVA_DB.Model(&mall.MallUser{})
	var mallUsers []mall.MallUser
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&mallUsers).Error
	return err, mallUsers, total
}

// LockUser 超级管理员修改管理员用户状态
func (m *ManageUserService) LockUser(idReq request.IdsReq, lockStatus int) (err error) {
	// 0 正常,1 禁止
	if lockStatus != 0 && lockStatus != 1 {
		return errors.New("操作非法！")
	}
	// 更新 用户表 UpdateColumns locked_flag
	err = global.GVA_DB.Model(&mall.MallUser{}).Where("user_id in ?", idReq.Ids).Update("locked_flag", lockStatus).Error
	return err
}
