/*
根据雪花算法生成唯一标识uuid
*/
package utils

import (
	"fmt"

	snowflake "github.com/bingjiekang/SnowFlake"
)

// 生成唯一标识的雪花id
func SnowFlakeUUid() int64 {
	// initialization
	snowf, err := snowflake.GetSnowFlake(0, "", "")
	if err != nil {
		fmt.Println(err)
	}
	// output ID
	return snowf.Generate()
}
