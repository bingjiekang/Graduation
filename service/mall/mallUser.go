package mall

import (
	_ "net/http"

	"github.com/gin-gonic/gin"
)

type MallUserService struct {
}

// 处理用户注册的数据库操作
func (m *MallUserService) Register(c *gin.Context) {

}

// 处理用户登陆的数据库操作
func (m *MallUserService) Login(c *gin.Context) {

}
