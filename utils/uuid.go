/*
根据雪花算法生成唯一标识uuid
*/
package utils

import (
	snowflake "github.com/bingjiekang/SnowFlake"
)

// initialization
var snowf, _ = snowflake.GetSnowFlake(0, "", "")

// 生成唯一标识的雪花id
func SnowFlakeUUid() int64 {
	// output ID
	return snowf.Generate()
}
