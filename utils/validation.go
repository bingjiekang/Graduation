package utils

import "regexp"

// 验证手机号是否合法
func ValidatePhoneNumber(phone string) bool {
	// 定义手机号格式的正则表达式
	pattern := `^1[3456789]\d{9}$`
	// 创建正则表达式对象并编译
	reg := regexp.MustCompile(pattern)
	// 判断手机号是否符合正则表达式
	if reg.MatchString(phone) {
		return true
	} else {
		return false
	}
}

// 验证密码是否符合要求(8位及以上,包含大小写字母或者数字和特殊字符)
func ValidatePassword(password string) bool {
	// 定义密码要求的正则表达式
	pattern := "^[A-Za-z\\d@$!%*#?&]{8,}$"

	// 创建正则表达式对象
	regExp := regexp.MustCompile(pattern)

	// 判断密码是否与正则表达式匹配
	if regExp.MatchString(password) {
		return true
	} else {
		return false
	}
}
