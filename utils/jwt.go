package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var signingKey = []byte("BlockChainMall")

// 自定义结构体
type Claim struct {
	Uuid     int64
	Jwtclaim jwt.StandardClaims
}

// 完成实现接口的Valid函数即可使用jwt.NewWithClaims
func (*Claim) Valid() error {
	return nil
}

// jwt加密
func CreateToken(uuid int64) (string, error) {
	// 使用结构体初始化信息
	claim := Claim{
		Uuid: uuid,
		Jwtclaim: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60,    // 1分钟前开始生效
			ExpiresAt: time.Now().Unix() + 60*60, // 1个小时后过期
			Issuer:    "AuthorJay",
		},
	}

	// SigningMethodHS256,HS256对称加密方式
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claim)
	// 通过自定义令牌加密
	key, err := token.SignedString(signingKey)
	if err != nil {
		fmt.Println("生成token失败")
	}
	return key, err

}

// jwt解密
func UndoToken(token string) (uuid int64, err error, ok bool) {
	Token, err := jwt.ParseWithClaims(token, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return 0, err, false
	}
	// 已经超时
	if time.Now().Unix() > Token.Claims.(*Claim).Jwtclaim.ExpiresAt {
		// fmt.Println("Token 已经超时!")
		return 0, fmt.Errorf("Token已超时!"), false
	}
	// 返回唯一标识Guid和管理员id
	return Token.Claims.(*Claim).Uuid, nil, true
}
