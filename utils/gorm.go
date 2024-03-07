package utils

import (
	"Graduation/global"
	internalmysql "Graduation/initialize/internalMysql"
	"Graduation/model/mall"
	"Graduation/model/manage"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Gorm 初始化数据库并产生数据库全局变量
func Gorm() *gorm.DB {
	switch global.GVA_CONFIG.System.DbType {
	case "mysql":
		return InitGormMysql()
	default:
		return InitGormMysql()
	}
}

// 链接mysql数据库
func InitGormMysql() *gorm.DB {
	m := global.GVA_CONFIG.Mysql
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         211,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	// 建立并打开对应的mysql文件柄
	if db, err := gorm.Open(mysql.New(mysqlConfig), internalmysql.Gorm.Config()); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		// 创建数据表...(many)
		{
			db.AutoMigrate(&mall.MallUser{})            // 用户信息表
			db.AutoMigrate(&mall.MallUserAddress{})     // 用户地址表
			db.AutoMigrate(&manage.MallCarousel{})      // 轮播图表
			db.AutoMigrate(&manage.MallIndexConfig{})   // 首页信息表
			db.AutoMigrate(&manage.MallGoodsCategory{}) // 分类信息表
			db.AutoMigrate(&manage.MallGoodsInfo{})     // 商品信息表
			db.AutoMigrate(&mall.MallShopCartItem{})    // 购物车信息表
			db.AutoMigrate(&manage.MallOrder{})         // 用户订单详情表
			db.AutoMigrate(&manage.MallOrderItem{})     // 用户订单列表
		}
		// 后台管理系统
		{
			db.AutoMigrate(&manage.MallAdminUser{}) // 管理员和超级管理员信息表
		}
		// db.AutoMigrate(&users.UserTrade{})
		global.GVA_LOG.Info("数据库连接成功!")
		return db
	}
}
