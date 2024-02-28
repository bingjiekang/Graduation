package mall

import (
	"Graduation/global"
	"Graduation/model/mall"
	"Graduation/model/mall/request"
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
func (m *MallUserService) RegisterUser(req request.RegisterUserParam) error {
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
func (m *MallUserService) LoginUser(req request.UserLoginParam) (err error, user mall.MallUser, token string) {
	err = global.GVA_DB.Where("login_name=? AND password_md5=?", req.LoginName, req.PasswordMd5).First(&user).Error

	// 生成对应token
	token, _ = utils.CreateToken(user.UUid)
	strUuid := fmt.Sprintf("%d", user.UUid)
	err = global.GVA_REDIS.Set(global.GVA_CTX, strUuid, token, 3600*time.Second).Err()
	if err != nil {
		global.GVA_LOG.Error("redis存储token失败")
	} else {
		global.GVA_LOG.Info("redis存储token成功")
	}

	return err, user, token
}

// 删除用户登陆token
func (m *MallUserService) DeleteMallUserToken(token string) (err error) {
	uuid, _, _ := utils.UndoToken(token)
	uid := fmt.Sprintf("%d", uuid)
	_, err = global.GVA_REDIS.Get(global.GVA_CTX, uid).Result()
	if err != nil {
		global.GVA_LOG.Error("无法从redis得到对应uid的token信息,可能已经失效或不存在")
		return err
	}
	// 存在则删除
	global.GVA_REDIS.Del(global.GVA_CTX, uid)
	return err
}
