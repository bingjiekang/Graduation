package mall

import (
	"Graduation/model/common/response"
	"fmt"

	"github.com/gin-gonic/gin"
)

type MallUserApi struct {
}

// 处理用户注册的路由接口转接
func (m *MallUserApi) UserRegister(c *gin.Context) {
	fmt.Println("接收到注册")
	// 检查用户传入的信息是否合理
	response.OkWithMessage("创建成功!", c)
	// 将数据加入到数据库
}

// 处理用户登录的路由接口转接
func (m *MallUserApi) UserLogin(c *gin.Context) {
	fmt.Println("接收登陆")
	response.OkWithMessage("登陆成功!", c)

}

// 处理用户退出的路由接口转接
func (m *MallUserApi) UserLogout(c *gin.Context) {
	token := c.GetHeader("token")
	fmt.Println("登陆token", token)
	// if err := mallUserTokenService.DeleteMallUserToken(token); err != nil {
	// 	response.FailWithMessage("登出失败", c)
	// } else {
	// 	response.OkWithMessage("登出成功", c)
	// }
	response.OkWithMessage("登出成功", c)

}
