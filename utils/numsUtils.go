package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// NumsInList 数值是否存在
func NumsInList(num int, nums []int) bool {
	for _, s := range nums {
		if s == num {
			return true
		}
	}
	return false
}

// GenValidateCode 随机6位数
func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, err := fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
		if err != nil {
			return ""
		}
	}
	return sb.String()
}

// 生成订单号
func GenOrderNo() string {
	// 获取当前日期时间作为订单号的一部分
	currentTime := time.Now().Format("20060102150405")

	// 生成随机数作为订单号的一部分
	rand.Seed(time.Now().UnixNano())
	randomPart := fmt.Sprintf("%04d", rand.Intn(10000))

	// 构建订单号
	orderNumber := currentTime + randomPart

	return orderNumber
}

// '2,3' 转换为[2,3]
func StrToList(strNum string) (nums []int) {
	strNums := strings.Split(strNum, ",")
	for _, s := range strNums {
		i, _ := strconv.Atoi(s)
		nums = append(nums, i)
	}
	return
}
