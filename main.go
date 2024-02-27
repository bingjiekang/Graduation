package main

import (
	"Graduation/global"
	"Graduation/initialize"
	"Graduation/utils"
)

func main() {
	// router := gin.Default()
	router := initialize.Routers()
	global.GVA_VIP = utils.Viper()        // 初始化Viper
	global.GVA_LOG = utils.Zap()          // 初始化zap日志库
	global.GVA_DB = utils.InitGormMysql() // gorm连接数据库
	router.Run(":8080")
}
