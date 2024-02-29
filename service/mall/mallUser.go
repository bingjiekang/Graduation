package mall

import (
	"Graduation/global"
	"Graduation/model/mall"
	requ "Graduation/model/mall/request"
	resp "Graduation/model/mall/response"
	"Graduation/utils"
	"errors"
	"fmt"
	_ "net/http"
	"time"

	"gorm.io/gorm"
)

type MallUserService struct {
}

// 处理用户注册的数据库操作
func (m *MallUserService) RegisterUser(req requ.RegisterUserParam) error {
	// 重复注册
	if !errors.Is(global.GVA_DB.Where("login_name =?", req.LoginName).First(&mall.MallUser{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同用户名")
	}
	// 注册成功
	return global.GVA_DB.Create(&mall.MallUser{
		UUid:          utils.SnowFlakeUUid(),
		LoginName:     req.LoginName,
		PasswordMd5:   utils.Md5(req.Password),
		IntroduceSign: "生命在于感知,生活在于选择!",
		// Create:    common.JSONTime{Time: time.Now()},
	}).Error
}

// 处理用户登陆的数据库操作
func (m *MallUserService) LoginUser(req requ.UserLoginParam) (err error, user mall.MallUser, token string) {
	err = global.GVA_DB.Where("login_name=? AND password_md5=?", req.LoginName, req.PasswordMd5).First(&user).Error
	if err != nil { // 没有找到登录信息
		return err, user, token
	}
	// 生成对应token
	token, _ = utils.CreateToken(user.UUid)
	strUuid := fmt.Sprintf("%d", user.UUid)
	err = global.GVA_REDIS.Set(global.GVA_CTX, strUuid, token, 3600*time.Second).Err()
	if err != nil {
		global.GVA_LOG.Error("redis存储token失败")
		return err, user, token
	} else {
		global.GVA_LOG.Info("redis存储token成功")
	}

	return err, user, token
}

// 删除用户登陆token
func (m *MallUserService) DeleteMallUserToken(token string) (err error) {
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

// 检查token是否存在
func (m *MallUserService) ExistUserToken(token string) (err error, tm int64) {
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

// 获取用户信息
func (m *MallUserService) GetUserInfo(token string) (err error, userInfoDetail resp.MallUserDetailResponse) {
	// 判断用户是否存在
	if !m.IsUserExist(token) {
		return errors.New("不存在的用户"), userInfoDetail
	}
	var userInfo mall.MallUser
	uuid, _, _ := utils.UndoToken(token)
	if err = global.GVA_DB.Where("u_uid =?", uuid).First(&userInfo).Error; err != nil {
		return errors.New("用户信息获取失败"), userInfoDetail
	}
	// 对应信息进行赋值
	{
		userInfoDetail.LoginName = userInfo.LoginName
		userInfoDetail.NickName = userInfo.NickName
		userInfoDetail.UUid = userInfo.UUid
		userInfoDetail.IntroduceSign = userInfo.IntroduceSign
	}
	return
}

// 更改用户信息
func (m *MallUserService) UpdateUserInfo(token string, req requ.UpdateUserInfoParam) (err error) {
	// 判断用户是否存在
	if !m.IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	var userInfo mall.MallUser
	uuid, _, _ := utils.UndoToken(token)
	if err = global.GVA_DB.Where("u_uid =?", uuid).First(&userInfo).Error; err != nil {
		return errors.New("用户信息获取失败(更改信息)")
	}
	// 若密码不为空，则表明用户修改密码
	{
		userInfo.NickName = req.NickName
		userInfo.IntroduceSign = req.IntroduceSign
	}
	if !(req.PasswordMd5 == "") {
		userInfo.PasswordMd5 = req.PasswordMd5
	}
	err = global.GVA_DB.Save(&userInfo).Error
	return
}

// 判断用户是否存在
func (m *MallUserService) IsUserExist(token string) bool {
	var userInfo mall.MallUser
	uuid, _, _ := utils.UndoToken(token)
	// uid := fmt.Sprintf("%d", uuid)
	if err := global.GVA_DB.Where("u_uid=?", uuid).First(&userInfo).Error; err != nil {
		return false // 用户不存在
	}
	return true
}
