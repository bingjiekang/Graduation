package internalmysql

import (
	"Graduation/global"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Gorm = new(_gorm)

type _gorm struct{}

// Config gorm 自定义配置
func (g *_gorm) Config() *gorm.Config {
	config := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NowFunc: func() time.Time {
			localSH, _ := time.LoadLocation("Asia/Shanghai")
			return time.Now().In(localSH).Local()
		},
	}
	_default := logger.New(NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})
	switch global.GVA_CONFIG.Mysql.LogMode {
	case "silent", "Silent":
		config.Logger = _default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = _default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = _default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = _default.LogMode(logger.Info)
	default:
		config.Logger = _default.LogMode(logger.Info)
	}
	return config
}
