package manage

import "time"

// MallAdminUser 结构体(包含管理员用户和超级管理员)
type MallAdminUser struct {
	Id            int       `json:"id" form:"id" gorm:"primarykey;AUTO_INCREMENT"`
	UUid          int64     `json:"uUid" form:"uUid" gorm:"column:u_uid;comment:唯一标识id"`
	LoginUserName string    `json:"loginUserName" form:"loginUserName" gorm:"column:login_user_name;comment:管理员登陆名称;type:varchar(50);"`
	LoginPassword string    `json:"loginPassword" form:"loginPassword" gorm:"column:login_password;comment:管理员登陆密码;type:varchar(50);"`
	NickName      string    `json:"nickName" form:"nickName" gorm:"column:nick_name;comment:管理员显示昵称;type:varchar(50);"`
	Locked        int       `json:"locked" form:"locked" gorm:"column:locked;comment:是否锁定 0未锁定 1已锁定无法登陆;type:tinyint"`
	IsSuperAdmin  int       `json:"isSuperAdmin" form:"isSuperAdmin" gorm:"column:is_super_admin;comment:是否为超级管理员 0普通管理员 1超级管理员;type:tinyint"`
	CreatedAt     time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;type:datetime"`
	UpdatedAt     time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:更新时间;type:datetime"`
}

func (MallAdminUser) TableName() string {
	return "super_admin_user"
}
