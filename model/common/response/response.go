/*
实现消息模版对不同操作的返回
*/
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	ResultCode int         `json:"resultCode"`
	Data       interface{} `json:"data"`
	Msg        string      `json:"message"`
}

const (
	ERROR   = 500
	SUCCESS = 200
	UNLOGIN = 416
)

// 返回结构的模版(返回对应信息)
func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		ResultCode: code,
		Data:       data,
		Msg:        msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "SUCCESS", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

func UnLogin(data interface{}, c *gin.Context) {
	Result(UNLOGIN, data, "未登录！", c)
}
