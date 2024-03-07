package request

import (
	"Graduation/model/common/request"
	"Graduation/model/mall"
)

type ManageLoginParam struct {
	UserName    string `json:"userName"`
	PasswordMd5 string `json:"passwordMd5"`
}
type ManageParam struct {
	LoginUserName string `json:"loginUserName"`
	LoginPassword string `json:"loginPassword"`
	NickName      string `json:"nickName"`
}

type ManageUpdateNameParam struct {
	NickName string `json:"nickName"`
}

type ManageUpdatePasswordParam struct {
	OriginalPassword string `json:"originalPassword"`
	NewPassword      string `json:"newPassword"`
}

type MallUserSearch struct {
	mall.MallUser
	request.PageInfo
}
