package global

import (
	"Graduation/config"

	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	// 默认配置文件
	ConfigFile = "config.yaml"
)

var (
	// 配置文件对应(类似指针,用来操作和读取对应信息)
	GVA_CONFIG config.Server
	// 全局 viper 用来指向对应信息
	GVA_VIP *viper.Viper
	// 全局 zap 用来全局记录日志
	GVA_LOG *zap.Logger
	// 全局 gorm 操作数据库
	GVA_DB *gorm.DB
	// 全局 redis 操作缓存数据库
	GVA_REDIS *redis.Client
	// 全局ctx
	GVA_CTX = context.Background()
)
