2024.2.26 开始毕业论文设计

# 代码的实现

## 1.安装Go环境

项目使用 Golang 1.22版本

## 2.下载并配置环境

``` 
下载地址：https://golang.google.cn/dl/
```	

### 2.1.确定使用Gin框架

```golang
	// 下载GIn框架包
	go get -u github.com/gin-gonic/gin
```

## 3.确定项目目录



## 4.代码编写

### 4.1.启动跨域请求

```golang
// 如果有跨域则打开
// Router = gin.Default()
// Router.Use(CrossDomain())

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 处理跨域请求,支持options访问
func CrossDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id,X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
```

### 4.2.完成viper和zap的使用，以及Gorm操作mysql的配置

```golang
	// viper 下载（便捷操作数据库和读写配置文件）
	go get github.com/spf13/viper
	// zap 下载 （日志文件系统）
	go get -u go.uber.org/zap
	// 用来分割日志
	go get -u github.com/natefinch/lumberjack
	// gorm 下载（操作数据库）
	go get -u gorm.io/gorm
	// 下载mysql和gorm相关
	go get -u gorm.io/driver/mysql
	// 配置viper
	如下 4.2.1
	
	// 配置Zap
	如下 4.2.2
	
	// 配置Gorm
	如下 4.2.3
	
	// 在main函数中启动并初始化他们三个

```

#### 4.2.1.配置 viper

```golang
package utils

import (
	"Graduation/global"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 用来配置并操作配置文件
func Viper(path ...string) *viper.Viper {
	var config string
	if len(path) == 0 { // 使用默认配置config.yaml
		fmt.Printf("您将使用默认的配置文件:[%v]\n", global.ConfigFile)
		config = global.ConfigFile
	} else { // 使用传入地址的配置文件
		fmt.Printf("您将使用指定配置文件:[%v]\n", path[0])
		config = path[0]
	}

	v := viper.New()
	// 设置配置文件
	v.SetConfigFile(config)
	// 覆盖和加载配置文件
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 监听配置 如果配置文件发送改变则变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
		fmt.Println(err)
	}
	return v

}

```

#### 4.2.2.配置 zap

```golang
package config

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`                           // 级别
	Format        string `mapstructure:"format" json:"format" yaml:"format"`                        // 输出
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                        // 日志前缀
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`                 // 日志文件夹
	ShowLine      bool   `mapstructure:"show-line" json:"showLine" yaml:"show-line"`                // 显示行
	EncodeLevel   string `mapstructure:"encode-level" json:"encodeLevel" yaml:"encode-level"`       // 编码级
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktraceKey" yaml:"stacktrace-key"` // 栈名
	LogInConsole  bool   `mapstructure:"log-in-console" json:"logInConsole" yaml:"log-in-console"`  // 输出控制台
}

```

```golang
package utils

import (
	"Graduation/global"
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 日志的操作和配置
func Zap() (logger *zap.Logger) {
	if ok, _ := PathExists(global.GVA_CONFIG.Zap.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("Create %v Directory\n", global.GVA_CONFIG.Zap.Director)
		err := os.Mkdir(global.GVA_CONFIG.Zap.Director, os.ModePerm)
		if err != nil {
			fmt.Println("创建文件夹失败!", err.Error())
		}
	}

	// 调试级别
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})
	// 日志级别
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.InfoLevel
	})
	// 警告级别
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.WarnLevel
	})
	// 错误级别
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	cores := [...]zapcore.Core{
		getEncoderCore(fmt.Sprintf("./%s/mallDebug.log", global.GVA_CONFIG.Zap.Director), debugPriority),
		getEncoderCore(fmt.Sprintf("./%s/mallInfo.log", global.GVA_CONFIG.Zap.Director), infoPriority),
		getEncoderCore(fmt.Sprintf("./%s/mallWarn.log", global.GVA_CONFIG.Zap.Director), warnPriority),
		getEncoderCore(fmt.Sprintf("./%s/mallError.log", global.GVA_CONFIG.Zap.Director), errorPriority),
	}
	logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())

	if global.GVA_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger

}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  global.GVA_CONFIG.Zap.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case global.GVA_CONFIG.Zap.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case global.GVA_CONFIG.Zap.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case global.GVA_CONFIG.Zap.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case global.GVA_CONFIG.Zap.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	if global.GVA_CONFIG.Zap.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer := GetWriteSyncer(fileName) // 使用file-rotatelogs进行日志分割
	return zapcore.NewCore(getEncoder(), writer, level)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	localSH, _ := time.LoadLocation("Asia/Shanghai")
	enc.AppendString(t.In(localSH).Format(global.GVA_CONFIG.Zap.Prefix + "2006/01/02 15:04:05"))
}

// @function: PathExists
// @description: 文件目录是否存在
// @param: path string
// @return: bool, error
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 分割日志
func GetWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file, //日志文件的位置
		MaxSize:    50,   //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 200,  //保留旧文件的最大个数
		MaxAge:     90,   //保留旧文件的最大天数
		Compress:   true, //是否压缩/归档旧文件
	}

	if global.GVA_CONFIG.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}
```

#### 4.2.3.配置 Gorm

```golang
package config

type Mysql struct {
	Path         string `mapstructure:"path" json:"path" yaml:"path"`                             // 服务器地址
	Port         string `mapstructure:"port" json:"port" yaml:"port"`                             // 端口
	Config       string `mapstructure:"config" json:"config" yaml:"config"`                       // 高级配置
	Dbname       string `mapstructure:"db-name" json:"dbname" yaml:"db-name"`                     // 数据库名
	Username     string `mapstructure:"username" json:"username" yaml:"username"`                 // 数据库用户名
	Password     string `mapstructure:"password" json:"password" yaml:"password"`                 // 数据库密码
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
	LogMode      string `mapstructure:"log-mode" json:"logMode" yaml:"log-mode"`                  // 是否开启Gorm全局日志
	LogZap       bool   `mapstructure:"log-zap" json:"logZap" yaml:"log-zap"`                     // 是否通过zap写入日志文件
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}

```

```golang
package utils

import (
	"Graduation/global"
	internalmysql "Graduation/initialize/internalMysql"
	"Graduation/model/mall"

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
		db.AutoMigrate(&mall.MallUser{}) // 用户信息表
		// db.AutoMigrate(&users.UserTrade{})
		global.GVA_LOG.Info("数据库连接成功!")
		return db
	}
}


```

### 4.3.登陆、注册、注销页面编写

#### jwt验证

```golang
// go get -u github.com/dgrijalva/jwt-go
package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var signingKey = []byte("BlockChainMall")

// 自定义结构体
type Claim struct {
	Uuid     int64
	Jwtclaim jwt.StandardClaims
}

// 完成实现接口的Valid函数即可使用jwt.NewWithClaims
func (*Claim) Valid() error {
	return nil
}

// jwt加密
func CreateToken(uuid int64) (string, error) {
	// 使用结构体初始化信息
	claim := Claim{
		Uuid: uuid,
		Jwtclaim: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60,    // 1分钟前开始生效
			ExpiresAt: time.Now().Unix() + 60*60, // 1小时后过期
			Issuer:    "AuthorJay",
		},
	}

	// SigningMethodHS256,HS256对称加密方式
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claim)
	// 通过自定义令牌加密
	key, err := token.SignedString(signingKey)
	if err != nil {
		fmt.Println("生成token失败")
	}
	return key, err

}

// jwt解密
func UndoToken(token string) (uuid int64, err error) {
	Token, err := jwt.ParseWithClaims(token, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	// 返回唯一标识Guid和管理员id
	return Token.Claims.(*Claim).Uuid, nil
}

```

#### uuid生成

```golang
//  go get "github.com/bingjiekang/SnowFlake"
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


```






