package utils

import "strconv"

// 任何类型转整形
func Transfer(data interface{}) (num int) {
	switch v := data.(type) {
	case float64:
		return int(v)
	case string:
		// 如果是字符串，直接使用
		num, _ = strconv.Atoi(v)
		return
	case int:
		return v
	}
	return
}
