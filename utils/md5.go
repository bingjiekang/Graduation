package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// 对密码进行加密的md5算法
func Md5(message string) string {
	// 创建一个新的hash对象将字符串转为字节切片
	hash := md5.Sum([]byte(message))
	// 将字节切片转为16进制字符串标识
	return hex.EncodeToString(hash[:])
}
