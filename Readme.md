2024.2.26 开始毕业论文设计

	暂定:区块链商城项目
	前端端口:xxx:8088
	后端端口:xxx:8080
	

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
	global.GVA_VIP = utils.Viper()        // 初始化Viper
	global.GVA_LOG = utils.Zap()          // 初始化zap日志库
	global.GVA_DB = utils.InitGormMysql() // gorm连接数据库
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

#### 4.3.1.jwt验证

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
			ExpiresAt: time.Now().Unix() + 60*60, // 1个小时后过期
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
func UndoToken(token string) (uuid int64, err error, ok bool) {
	Token, err := jwt.ParseWithClaims(token, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return 0, err, false
	}
	// 已经超时
	if time.Now().Unix() > Token.Claims.(*Claim).Jwtclaim.ExpiresAt {
		// fmt.Println("Token 已经超时!")
		return 0, fmt.Errorf("Token已超时!"), false
	}
	// 返回唯一标识Guid和管理员id
	return Token.Claims.(*Claim).Uuid, nil, true
}

```

#### 4.3.2.uuid生成

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

// initialization
var snowf, _ = snowflake.GetSnowFlake(0, "", "")

// 生成唯一标识的雪花id
func SnowFlakeUUid() int64 {
	// output ID
	return snowf.Generate()
}


```

#### 4.3.3.redis 配置和初始化

```golang
// 下载 redis 依赖包
//  go get github.com/redis/go-redis/v9

package utils

import (
	"Graduation/global"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func RedisConnect() *redis.Client {
	var r = global.GVA_CONFIG.Redis
	options := redis.Options{
		Addr:            r.RedisHost + ":" + r.RedisPort,
		DB:              r.RedisDb,
		PoolSize:        r.RedisPoolSize,                                 // Redis连接池大小
		MaxRetries:      r.RedisMaxRetries,                               // 最大重试次数
		ConnMaxIdleTime: time.Second * time.Duration(r.RedisIdleTimeout), // 空闲链接超时时间
	}
	if r.RedisPassword != "" {
		options.Password = r.RedisPassword
	}
	Rdb := redis.NewClient(&options)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pong, err := Rdb.Ping(ctx).Result()
	if err == redis.Nil {
		global.GVA_LOG.Debug("[StoreRedis] Nil reply returned by Rdb when key does not exist.")
	} else if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("[StoreRedis] redis connRdb err,err=%s", err))
		panic(err)
	} else {
		global.GVA_LOG.Debug(fmt.Sprintf("[StoreRedis] redis connRdb success,suc=%s", pong))
	}
	return Rdb

}

```

```golang

package config
// 缓存数据库
type Redis struct {
	RedisHost        string `mapstructure:"redis-host" json:"redisHost" yaml:"redis-host"` // redis地址
	RedisPort        string `mapstructure:"redis-port" json:"redisPort" yaml:"redis-port"`
	RedisPassword    string `mapstructure:"redis-password" json:"redisPassword" yaml:"redis-password"`
	RedisDb          int    `mapstructure:"redis-db" json:"redisDb" yaml:"redis-db"`
	RedisPoolSize    int    `mapstructure:"redis-pool-size" json:"redisPoolSize" yaml:"redis-pool-size"`
	RedisMaxRetries  int    `mapstructure:"redis-max-retries" json:"redisMaxRetries" yaml:"redis-max-retries"`
	RedisIdleTimeout int    `mapstructure:"redis-idle-timeout" json:"redisIdleTimeout" yaml:"redis-idle-timeout"`
}

```

#### 4.3.4.手机号以及密码验证

```golang

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
	
```

#### 4.3.5md5 加密

```golang
// 对密码进行加密的md5算法
func Md5(message string) string {
	// 创建一个新的hash对象将字符串转为字节切片
	hash := md5.Sum([]byte(message))
	// 将字节切片转为16进制字符串标识
	return hex.EncodeToString(hash[:])
}
```


#### 4.3.6.登陆、注册、登出代码接口编写

```golang

package mall

import (
	"Graduation/global"
	"Graduation/model/common/response"
	"Graduation/model/mall/request"
	"Graduation/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MallUserApi struct {
}

// 处理用户注册的路由接口转接
func (m *MallUserApi) UserRegister(c *gin.Context) {
	// 绑定对应用户注册信息结构体
	var req request.RegisterUserParam
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("注册信息无法绑定对应结构")
	}
	// 检查用户传入的信息是否合理
	if !(utils.ValidatePhoneNumber(req.LoginName) && utils.ValidatePassword(req.Password)) { // 有一个不合理
		response.FailWithMessage("请确保用户名为手机号,密码为8位以上数字,密码,特殊符合的组合!", c)
		return
	}
	// 对用户进行检查和注册处理
	if err := mallUserService.RegisterUser(req); err != nil {
		global.GVA_LOG.Error("创建失败", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
	}
	response.OkWithMessage("创建成功!", c)
}

// 处理用户登录的路由接口转接
func (m *MallUserApi) UserLogin(c *gin.Context) {
	// 绑定对应登录信息结构体
	var req request.UserLoginParam
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("登陆信息无法绑定对应结构")
	}
	// 校验登陆信息是否正确
	if err, _, token := mallUserService.LoginUser(req); err != nil {
		response.FailWithMessage("登陆失败!请检查账号和密码是否错误", c)
	} else {
		response.OkWithData(token, c)
	}
}

// 处理用户退出的路由接口转接
func (m *MallUserApi) UserLogout(c *gin.Context) {
	token := c.GetHeader("token")
	global.GVA_LOG.Info("登陆token:" + token)
	// 检查并删除用户 token
	if err := mallUserService.DeleteMallUserToken(token); err != nil {
		response.FailWithMessage("登出失败", c)
	} else {
		response.OkWithMessage("登出成功", c)
	}
}

// 处理用户信息的路由接口转接
func (m *MallUserApi) UserInfo(c *gin.Context) {
	token := c.GetHeader("token")
	// 获取用户信息并返回
	if err, userDetail := mallUserService.GetUserInfo(token); err != nil {
		global.GVA_LOG.Error("未查询到用户记录", zap.Error(err))
		response.FailWithMessage("未查询到用户记录", c)
	} else {
		response.OkWithData(userDetail, c)
	}
}

// 用户修改信息接口路由
func (m *MallUserApi) UpdateUserInfo(c *gin.Context) {
	token := c.GetHeader("token")
	var req request.UpdateUserInfoParam
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("用户更新的信息无法绑定对应结构")
	}
	// 获取用户信息
	if err := mallUserService.UpdateUserInfo(token, req); err != nil {
		global.GVA_LOG.Error("更新用户信息失败", zap.Error(err))
		response.FailWithMessage("更新用户信息失败"+err.Error(), c)
	}
	response.OkWithMessage("更新成功", c)
}
```

#### 4.3.7.登陆、注册、登出代码相关数据库操作编写

```golang
package mall

import (
	"Graduation/global"
	"Graduation/model/mall"
	requ "Graduation/model/mall/request"
	resp "Graduation/model/mall/response"
	"Graduation/utils"
	"errors"
	"fmt"
	_ "net/http"
	"time"

	"gorm.io/gorm"
)

type MallUserService struct {
}

// 处理用户注册的数据库操作
func (m *MallUserService) RegisterUser(req requ.RegisterUserParam) error {
	// 重复注册
	if !errors.Is(global.GVA_DB.Where("login_name =?", req.LoginName).First(&mall.MallUser{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同用户名")
	}
	// 注册成功
	return global.GVA_DB.Create(&mall.MallUser{
		UUid:          utils.SnowFlakeUUid(),
		LoginName:     req.LoginName,
		PasswordMd5:   utils.Md5(req.Password),
		IntroduceSign: "生命在于感知,生活在于选择!",
		// Create:    common.JSONTime{Time: time.Now()},
	}).Error
}

// 处理用户登陆的数据库操作
func (m *MallUserService) LoginUser(req requ.UserLoginParam) (err error, user mall.MallUser, token string) {
	err = global.GVA_DB.Where("login_name=? AND password_md5=?", req.LoginName, req.PasswordMd5).First(&user).Error
	if err != nil { // 没有找到登录信息
		return err, user, token
	}
	// 生成对应token
	token, _ = utils.CreateToken(user.UUid)
	strUuid := fmt.Sprintf("%d", user.UUid)
	err = global.GVA_REDIS.Set(global.GVA_CTX, strUuid, token, 3600*time.Second).Err()
	if err != nil {
		global.GVA_LOG.Error("redis存储token失败")
		return err, user, token
	} else {
		global.GVA_LOG.Info("redis存储token成功")
	}

	return err, user, token
}

// 删除用户登陆token
func (m *MallUserService) DeleteMallUserToken(token string) (err error) {
	uuid, err, ok := utils.UndoToken(token)
	if err != nil && ok == 0 { // 解码token出现错误
		return err
	}
	uid := fmt.Sprintf("%d", uuid)
	if ok == 1 { // 超时 token已经失效
		global.GVA_REDIS.Del(global.GVA_CTX, uid)
		return nil
	}
	_, err = global.GVA_REDIS.Get(global.GVA_CTX, uid).Result()
	if err != nil {
		global.GVA_LOG.Error("无法从redis得到对应uid的token信息,可能不存在")
		return err
	}
	// 存在则删除
	global.GVA_REDIS.Del(global.GVA_CTX, uid)
	return nil
}

// 检查token是否存在
func (m *MallUserService) ExistUserToken(token string) (err error, tm int64) {
	uuid, err, ok := utils.UndoToken(token)
	if err != nil && ok == 0 { // 解码token出现错误
		return err, 0
	}
	uid := fmt.Sprintf("%d", uuid)
	_, err = global.GVA_REDIS.Get(global.GVA_CTX, uid).Result()
	if err != nil {
		global.GVA_LOG.Error("token不存在")
		return err, 0
	}
	if ok == 1 {
		global.GVA_LOG.Info("token已超时")
		return err, 1
	}
	return nil, 2
}

// 获取用户信息
func (m *MallUserService) GetUserInfo(token string) (err error, userInfoDetail resp.MallUserDetailResponse) {
	// 判断用户是否存在
	if !m.IsUserExist(token) {
		return errors.New("不存在的用户"), userInfoDetail
	}
	var userInfo mall.MallUser
	uuid, _, _ := utils.UndoToken(token)
	if err = global.GVA_DB.Where("u_uid =?", uuid).First(&userInfo).Error; err != nil {
		return errors.New("用户信息获取失败"), userInfoDetail
	}
	// 对应信息进行赋值
	{
		userInfoDetail.LoginName = userInfo.LoginName
		userInfoDetail.NickName = userInfo.NickName
		userInfoDetail.UUid = userInfo.UUid
		userInfoDetail.IntroduceSign = userInfo.IntroduceSign
	}
	return
}

// 更改用户信息
func (m *MallUserService) UpdateUserInfo(token string, req requ.UpdateUserInfoParam) (err error) {
	// 判断用户是否存在
	if !m.IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	var userInfo mall.MallUser
	uuid, _, _ := utils.UndoToken(token)
	if err = global.GVA_DB.Where("u_uid =?", uuid).First(&userInfo).Error; err != nil {
		return errors.New("用户信息获取失败(更改信息)")
	}
	// 若密码不为空，则表明用户修改密码
	{
		userInfo.NickName = req.NickName
		userInfo.IntroduceSign = req.IntroduceSign
	}
	if !(req.PasswordMd5 == "") {
		userInfo.PasswordMd5 = req.PasswordMd5
	}
	err = global.GVA_DB.Save(&userInfo).Error
	return
}

// 判断用户是否存在
func (m *MallUserService) IsUserExist(token string) bool {
	var userInfo mall.MallUser
	uuid, _, _ := utils.UndoToken(token)
	// uid := fmt.Sprintf("%d", uuid)
	if err := global.GVA_DB.Where("u_uid=?", uuid).First(&userInfo).Error; err != nil {
		return false // 用户不存在
	}
	return true
}
```

### 4.4.用户地址管理代码编写

#### 4.4.1.增、删、改、查用户地址

```golang
// go get -u "github.com/jinzhu/copier" copier函数 复制对应相同的字段

package mall

import (
	"Graduation/global"
	"Graduation/model/common/response"
	"Graduation/model/mall/request"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MallUserAddressApi struct {
}

// 增加用户地址信息
func (m *MallUserAddressApi) AddUserAddress(c *gin.Context) {
	var req request.AddAddressParam
	_ = c.ShouldBindJSON(&req)
	token := c.GetHeader("token")
	// 保存用户地址信息
	err := mallUserAddressService.AddUserAddress(token, req)
	if err != nil {
		global.GVA_LOG.Error("创建失败", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// 查询用户地址列表信息
func (m *MallUserAddressApi) GetAddressList(c *gin.Context) {
	token := c.GetHeader("token")
	if err, userAddressList := mallUserAddressService.GetUserAddressList(token); err != nil {
		global.GVA_LOG.Error("获取地址列表信息失败", zap.Error(err))
		response.FailWithMessage("获取地址列表信息失败:"+err.Error(), c)
	} else if len(userAddressList) == 0 {
		global.GVA_LOG.Info("获取地址列表信息为空")
		response.OkWithData(nil, c)
	} else {
		response.OkWithData(userAddressList, c)
	}
}

// 查看用户指定标识的地址信息
func (m *MallUserAddressApi) GetUserAddress(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("addressId"), 10, 64)
	token := c.GetHeader("token")
	if err, userAddress := mallUserAddressService.GetUserAddress(token, id); err != nil {
		global.GVA_LOG.Error("获取指定地址信息失败", zap.Error(err))
		response.FailWithMessage("获取指定地址信息失败:"+err.Error(), c)
	} else {
		response.OkWithData(userAddress, c)
	}
}

// 修改用户地址信息
func (m *MallUserAddressApi) UpdateUserAddress(c *gin.Context) {
	// 接受修改的用户信息并绑定对应结构体
	var req request.UpdateAddressParam
	_ = c.ShouldBindJSON(&req)
	token := c.GetHeader("token")
	// 修改用户地址信息
	err := mallUserAddressService.UpdateUserAddress(token, req)
	if err != nil {
		global.GVA_LOG.Error("用户地址信息修改失败", zap.Error(err))
		response.FailWithMessage("用户地址信息修改失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("用户地址信息修改成功", c)
	}
}

// 删除指定用户标识的地址信息
func (m *MallUserAddressApi) DeleteUserAddress(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("addressId"), 10, 64)
	token := c.GetHeader("token")
	if err := mallUserAddressService.DeleteUserAddress(token, id); err != nil {
		global.GVA_LOG.Error("删除用户指定地址信息失败", zap.Error(err))
		response.FailWithMessage("删除用户指定地址信息失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("删除用户地址成功", c)
	}
}

```

#### 4.4.2.对应增删改查的数据库操作方法

```golang
package mall

import (
	"Graduation/global"
	"Graduation/model/mall"
	requ "Graduation/model/mall/request"
	"Graduation/utils"
	"errors"

	"github.com/jinzhu/copier"
)

type MallUserAddressService struct {
}

// 用户地址保存
func (m *MallUserAddressService) AddUserAddress(token string, req requ.AddAddressParam) (err error) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("用户不存在")
	}
	// 用户地址信息
	var userAddress mall.MallUserAddress
	err = copier.Copy(&userAddress, &req)
	if err != nil {
		return err
	}
	uuid, _, _ := utils.UndoToken(token)
	userAddress.Uuid = uuid
	// 判断是否为默认地址
	if req.DefaultFlag == 1 { // 新增默认地址
		// 查询是否已有默认地址
		if err = UpdateUserDefaultAddress(uuid); err != nil {
			return err
		}
		// // 查询是否已有默认地址
		// var defaultUserAddress mall.MallUserAddress
		// global.GVA_DB.Where("u_uid=? and default_flag =1 and is_deleted = 0", uuid).First(&defaultUserAddress)
		// // 已有默认地址(将原来默认地址取消)
		// if defaultUserAddress != (mall.MallUserAddress{}) {
		// 	defaultUserAddress.DefaultFlag = 0 // 设为非默认
		// 	err = global.GVA_DB.Save(&defaultUserAddress).Error
		// 	if err != nil {
		// 		return
		// 	}
		// }

	}
	// 创建新的地址
	if err = global.GVA_DB.Create(&userAddress).Error; err != nil {
		return
	}
	return
}

// GetUserAddressList 获取全部收货地址
func (m *MallUserAddressService) GetUserAddressList(token string) (err error, userAddress []mall.MallUserAddress) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("用户不存在"), userAddress
	}
	uuid, _, _ := utils.UndoToken(token)
	// 得到用户全部收货地址信息
	global.GVA_DB.Where("u_uid=? and is_deleted=0", uuid).Find(&userAddress)
	return
}

// 查询对应标识地址的信息
func (m *MallUserAddressService) GetUserAddress(token string, id int64) (err error, userAddress mall.MallUserAddress) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("用户不存在"), userAddress
	}
	uuid, _, _ := utils.UndoToken(token)
	// 得到用户对应标识id的收货地址信息
	global.GVA_DB.Where("u_uid = ? and address_id = ? and is_deleted = 0", uuid, id).Find(&userAddress)
	return
}

// 修改对应标识地址的信息
func (m *MallUserAddressService) UpdateUserAddress(token string, req requ.UpdateAddressParam) (err error) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("用户不存在")
	}
	// 解析 token 并获得uuid
	uuid, _, _ := utils.UndoToken(token)
	// 修改对应标识的 地址信息
	var reqUserAddr mall.MallUserAddress
	// 获取对应uid和对应标识的地址信息
	if err = global.GVA_DB.Where("address_id = ? and u_uid = ?", req.AddressId, uuid).First(&reqUserAddr).Error; err != nil {
		// 如果查询不到 证明传入的地址信息标识错误
		return errors.New("用户地址不存在")
	}
	// 成功获取到对应标识的地址信息
	if uuid != reqUserAddr.Uuid { // 不是用户本人
		return errors.New("非用户本人 禁止该操作！")
	}
	// 将对应修改的信息保存到 准备修改的这个对应标识的结构体里
	err = copier.Copy(&reqUserAddr, &req)
	if err != nil {
		return
	}
	// 修改对应数据信息 并保存
	if req.DefaultFlag == 1 { // 如果是默认收货地址
		// 查询是否已有默认地址
		if err = UpdateUserDefaultAddress(uuid); err != nil {
			return err
		}
	}
	// 将对应的uuid赋值进去
	reqUserAddr.Uuid = uuid
	err = global.GVA_DB.Save(&reqUserAddr).Error
	return
}

// 查询是否有默认地址信息 有默认地址信息则直接修改,
func UpdateUserDefaultAddress(uuid int64) (err error) {
	// 查询是否已有默认地址
	var defaultUserAddress mall.MallUserAddress
	global.GVA_DB.Where("u_uid=? and default_flag =1 and is_deleted = 0", uuid).First(&defaultUserAddress)
	// 已有默认地址(将原来默认地址取消)
	if defaultUserAddress != (mall.MallUserAddress{}) {
		defaultUserAddress.DefaultFlag = 0 // 设为非默认
		err = global.GVA_DB.Save(&defaultUserAddress).Error
		if err != nil {
			return
		}
	}
	return nil
}

// 删除对应标识的用户地址信息
func (m *MallUserAddressService) DeleteUserAddress(token string, id int64) (err error) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("用户不存在")
	}
	// 解析 token 并获得uuid
	uuid, _, _ := utils.UndoToken(token)
	// 修改对应标识的 地址信息
	var reqUserAddr mall.MallUserAddress
	// 获取对应uid和对应标识的地址信息
	if err = global.GVA_DB.Where("address_id = ? and u_uid = ?", id, uuid).First(&reqUserAddr).Error; err != nil {
		// 如果查询不到 证明传入的地址信息标识错误
		return errors.New("用户地址不存在")
	}
	// 成功获取到对应标识的地址信息
	if uuid != reqUserAddr.Uuid { // 不是用户本人
		return errors.New("非用户本人 禁止该操作！")
	}
	err = global.GVA_DB.Delete(&reqUserAddr).Error
	return
}

```

#### 4.4.3.对应修改地址信息的路由建立

```golang
package mall

import (
	v1 "Graduation/api/v1"
	"Graduation/middleware"

	"github.com/gin-gonic/gin"
)

// 用户地址路由连接表
type MallUserAddressRouter struct {
}

func (m *MallUserRouter) ApiMallUserAddressRouter(Router *gin.RouterGroup) {
	mallUserAddressRouter := Router.Group("v1").Use(middleware.UserJWTAuth())
	var userAddressApi = v1.ApiGroupApp.MallApiGroup.MallUserAddressApi
	{
		mallUserAddressRouter.POST("/address", userAddressApi.AddUserAddress)                 // 增加地址
		mallUserAddressRouter.GET("/address", userAddressApi.GetAddressList)                  // 查看用户全部地址列表信息
		mallUserAddressRouter.GET("/address/:addressId", userAddressApi.GetUserAddress)       // 获取指定地址详情
		mallUserAddressRouter.PUT("/address", userAddressApi.UpdateUserAddress)               // 修改用户指定地址信息
		mallUserAddressRouter.DELETE("/address/:addressId", userAddressApi.DeleteUserAddress) //删除地址
		// mallUserAddressRouter.GET("/address/default", userAddressApi.GetMallUserDefaultAddress) //获取默认地址

	}

}

```

### 4.5.首页信息显示

#### 4.5.1.配置首页对应信息路由


```golang

// 轮播商品展示
	err, _, mallCarouseInfo := mallCarouselService.GetIndexCarousels(5)
	if err != nil {
		global.GVA_LOG.Error("轮播图获取失败"+err.Error(), zap.Error(err))
		response.FailWithMessage("轮播图获取失败", c)
	}
	// 新品上线展示
	err, newGoodses := mallIndexConfigService.GetIndexInfomation(enum.IndexGoodsNew.Code(), 5)
	if err != nil {
		global.GVA_LOG.Error("新品获取失败"+err.Error(), zap.Error(err))
		response.FailWithMessage("新品获取失败", c)
	}
	// 热门商品展示
	err, hotGoodses := mallIndexConfigService.GetIndexInfomation(enum.IndexGoodsHot.Code(), 4)
	if err != nil {
		global.GVA_LOG.Error("热门商品获取失败"+err.Error(), zap.Error(err))
		response.FailWithMessage("热门商品获取失败", c)
	}
	// 最新推荐商品展示
	err, recommendGoodses := mallIndexConfigService.GetIndexInfomation(enum.IndexGoodsRecommond.Code(), 10)
	if err != nil {
		global.GVA_LOG.Error("推荐商品获取失败"+err.Error(), zap.Error(err))
		response.FailWithMessage("推荐商品获取失败", c)
	}
	// 首页全部商品数据
	indexResult := make(map[string]interface{})
	indexResult["carousels"] = mallCarouseInfo         // 轮播图数据
	indexResult["newGoodses"] = newGoodses             // 新品上市
	indexResult["hotGoodses"] = hotGoodses             // 热门商品
	indexResult["recommendGoodses"] = recommendGoodses // 推荐商品
	response.OkWithData(indexResult, c)

```


#### 4.5.2.首页信息对应数据库操作

```golang
package mall

import (
	"Graduation/global"
	"Graduation/model/mall"
	"Graduation/model/mall/response"
)

type MallCarouselService struct {
}

// GetIndexCarousels 首页返回固定数量的轮播图对象
func (m *MallCarouselService) GetIndexCarousels(num int) (err error, mallCarousels []mall.MallCarousel, list interface{}) {
	var carouselIndexs []response.MallIndexCarouselResponse
	err = global.GVA_DB.Where("is_deleted = 0").Order("carousel_rank desc").Limit(num).Find(&mallCarousels).Error
	for _, carousel := range mallCarousels {
		carouselIndex := response.MallIndexCarouselResponse{
			CarouselUrl: carousel.CarouselUrl,
			RedirectUrl: carousel.RedirectUrl,
		}
		carouselIndexs = append(carouselIndexs, carouselIndex)
	}
	return err, mallCarousels, carouselIndexs
}

```

```golang
package mall

import (
	"Graduation/global"
	"Graduation/model/mall"
	"Graduation/model/mall/response"
	"Graduation/utils"
)

type MallIndexInfomationService struct {
}

// GetIndexInfomation 首页新品/热门/推荐返回相关IndexConfig
func (m *MallIndexInfomationService) GetIndexInfomation(configType int, num int) (err error, list interface{}) {
	var indexConfigs []mall.MallIndexConfig
	err = global.GVA_DB.Where("config_type = ?", configType).Where("is_deleted = 0").Order("config_rank desc").Limit(num).Find(&indexConfigs).Error
	if err != nil {
		return
	}
	// 获取商品id
	var ids []int
	for _, indexConfig := range indexConfigs {
		ids = append(ids, indexConfig.GoodsId)
	}
	// 获取商品信息
	var goodsList []mall.MallGoodsInfo
	err = global.GVA_DB.Where("goods_id in ?", ids).Find(&goodsList).Error
	var indexGoodsList []response.MallIndexConfigGoodsResponse
	// 超出30个字符显示....
	for _, indexGoods := range goodsList {
		res := response.MallIndexConfigGoodsResponse{
			GoodsId:       indexGoods.GoodsId,
			GoodsName:     utils.ReplaceLength(indexGoods.GoodsName, 30),
			GoodsIntro:    utils.ReplaceLength(indexGoods.GoodsIntro, 30),
			GoodsCoverImg: indexGoods.GoodsCoverImg,
			SellingPrice:  indexGoods.SellingPrice,
			Tag:           indexGoods.Tag,
		}
		indexGoodsList = append(indexGoodsList, res)
	}
	return err, indexGoodsList
}
```

### 4.6.分类信息获取

#### 4.6.1.分类页信息转接路由

```golang
package mall

import (
	"Graduation/global"
	"Graduation/model/common/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MallGoodsCategoryApi struct {
}

// 返回分类数据 (分类页调用)
func (m *MallGoodsCategoryApi) GetGoodsCategorize(c *gin.Context) {
	err, list := mallGoodsCategoryService.GetGoodsCategories()
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败"+err.Error(), c)
	}
	response.OkWithData(list, c)
}

```

#### 4.6.2.分类页数据库操作

```golang
package mall

import (
	"Graduation/global"
	"Graduation/model/mall"
	"Graduation/model/mall/response"
	"Graduation/utils/enum"

	"github.com/jinzhu/copier"
)

type MallGoodsCategoryService struct {
}

// 获取分类页 一二三级 分类信息
func (m *MallGoodsCategoryService) GetGoodsCategories() (err error, MallIndexCategoryVOS []response.MallIndexCategoryVO) {

	// 获取并添加一级分类的固定数量的数据
	_, firstLevelCategories := selectByLevelAndParentIdsAndNumber([]int{0}, enum.LevelOne.Code(), 10)
	if firstLevelCategories != nil {
		var firstLevelCategoryIds []int
		for _, firstLevelCategory := range firstLevelCategories {
			firstLevelCategoryIds = append(firstLevelCategoryIds, firstLevelCategory.CategoryId)
		}
		// 获取并添加二级分类的数据
		_, secondLevelCategories := selectByLevelAndParentIdsAndNumber(firstLevelCategoryIds, enum.LevelTwo.Code(), 0)
		if secondLevelCategories != nil {
			var secondLevelCategoryIds []int
			for _, secondLevelCategory := range secondLevelCategories {
				secondLevelCategoryIds = append(secondLevelCategoryIds, secondLevelCategory.CategoryId)
			}
			// 获取并添加三级分类的数据
			_, thirdLevelCategories := selectByLevelAndParentIdsAndNumber(secondLevelCategoryIds, enum.LevelThree.Code(), 0)
			if thirdLevelCategories != nil {
				// 根据 parentId 将 thirdLevelCategories 分组
				thirdLevelCategoryMap := make(map[int][]mall.MallGoodsCategory)
				for _, thirdLevelCategory := range thirdLevelCategories {
					thirdLevelCategoryMap[thirdLevelCategory.ParentId] = []mall.MallGoodsCategory{}
				}
				for k, v := range thirdLevelCategoryMap {
					for _, third := range thirdLevelCategories {
						if k == third.ParentId {
							v = append(v, third)
						}
						thirdLevelCategoryMap[k] = v
					}
				}
				var secondLevelCategoryVOS []response.SecondLevelCategoryVO
				// 处理二级分类
				for _, secondLevelCategory := range secondLevelCategories {
					var secondLevelCategoryVO response.SecondLevelCategoryVO
					err = copier.Copy(&secondLevelCategoryVO, &secondLevelCategory)
					// 如果该二级分类下有数据则放入 secondLevelCategoryVOS 对象中
					if _, ok := thirdLevelCategoryMap[secondLevelCategory.CategoryId]; ok {
						// 根据二级分类的 id 取出 thirdLevelCategoryMap 分组中的三级分类 list
						tempGoodsCategories := thirdLevelCategoryMap[secondLevelCategory.CategoryId]
						var thirdLevelCategoryRes []response.ThirdLevelCategoryVO
						err = copier.Copy(&thirdLevelCategoryRes, &tempGoodsCategories)
						secondLevelCategoryVO.ThirdLevelCategoryVOS = thirdLevelCategoryRes
						secondLevelCategoryVOS = append(secondLevelCategoryVOS, secondLevelCategoryVO)
					}

				}
				//处理一级分类
				if secondLevelCategoryVOS != nil {
					//根据 parentId 将 thirdLevelCategories 分组
					secondLevelCategoryVOMap := make(map[int][]response.SecondLevelCategoryVO)
					for _, secondLevelCategory := range secondLevelCategoryVOS {
						secondLevelCategoryVOMap[secondLevelCategory.ParentId] = []response.SecondLevelCategoryVO{}
					}
					for k, v := range secondLevelCategoryVOMap {
						for _, second := range secondLevelCategoryVOS {
							if k == second.ParentId {
								var secondLevelCategory response.SecondLevelCategoryVO
								copier.Copy(&secondLevelCategory, &second)
								v = append(v, secondLevelCategory)
							}
							secondLevelCategoryVOMap[k] = v
						}
					}
					for _, firstCategory := range firstLevelCategories {
						var newBeeMallIndexCategoryVO response.MallIndexCategoryVO
						err = copier.Copy(&newBeeMallIndexCategoryVO, &firstCategory)
						//如果该一级分类下有数据则放入 MallIndexCategoryVOS 对象中
						if _, ok := secondLevelCategoryVOMap[firstCategory.CategoryId]; ok {
							//根据一级分类的id取出secondLevelCategoryVOMap分组中的二级级分类list
							tempGoodsCategories := secondLevelCategoryVOMap[firstCategory.CategoryId]
							newBeeMallIndexCategoryVO.SecondLevelCategoryVOS = tempGoodsCategories
							MallIndexCategoryVOS = append(MallIndexCategoryVOS, newBeeMallIndexCategoryVO)
						}
					}
				}
			}
		}
	}
	return
}

// 获取分类数据
func selectByLevelAndParentIdsAndNumber(ids []int, level int, limit int) (err error, categories []mall.MallGoodsCategory) {
	// 获取对应分类数据
	err = global.GVA_DB.Where("parent_id in ? and category_level =? and is_deleted = 0", ids, level).Order("category_rank desc").Limit(limit).Find(&categories).Error
	return
}

```

### 4.7.商品页面详细信息和商品搜索获取

#### 4.7.1.商业页面详情及搜索接口

```golang
package mall

import (
	"Graduation/global"
	"Graduation/model/common/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MallGoodsInfoApi struct {
}

// 商品详情页
func (m *MallGoodsInfoApi) GoodsDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err, goodsInfo := mallGoodsInfoService.GetMallGoodsInfo(id)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败"+err.Error(), c)
	}
	response.OkWithData(goodsInfo, c)
}

// 商品搜索
func (m *MallGoodsInfoApi) GoodsSearch(c *gin.Context) {
	// 获取页码
	pageNumber, _ := strconv.Atoi(c.Query("pageNumber"))
	goodsCategoryId, _ := strconv.Atoi(c.Query("goodsCategoryId"))
	keyword := c.Query("keyword")
	orderBy := c.Query("orderBy")
	if err, list, total := mallGoodsInfoService.MallGoodsListBySearch(pageNumber, goodsCategoryId, keyword, orderBy); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:       list,
			TotalCount: total,
			CurrPage:   pageNumber,
			PageSize:   10,
		}, "获取成功", c)
	}
}

```

#### 4.7.2.商品详情及搜索数据库操作

```golang
package mall

import (
	"Graduation/global"
	"Graduation/model/mall"
	"Graduation/model/mall/response"
	"Graduation/utils"
	"errors"

	"github.com/jinzhu/copier"
)

type MallGoodsInfoService struct {
}

// GetMallGoodsInfo 获取商品信息
func (m *MallGoodsInfoService) GetMallGoodsInfo(id int) (err error, res response.GoodsInfoDetailResponse) {
	var mallGoodsInfo mall.MallGoodsInfo
	err = global.GVA_DB.Where("goods_id = ?", id).First(&mallGoodsInfo).Error
	if mallGoodsInfo.GoodsSellStatus != 0 {
		return errors.New("商品已下架"), response.GoodsInfoDetailResponse{}
	}
	err = copier.Copy(&res, &mallGoodsInfo)
	if err != nil {
		return err, response.GoodsInfoDetailResponse{}
	}
	var list []string
	list = append(list, mallGoodsInfo.GoodsCarousel)
	res.GoodsCarouselList = list

	return
}

// MallGoodsListBySearch 商品搜索分页
func (m *MallGoodsInfoService) MallGoodsListBySearch(pageNumber int, goodsCategoryId int, keyword string, orderBy string) (err error, searchGoodsList []response.GoodsSearchResponse, total int64) {
	// 根据搜索条件查询
	var goodsList []mall.MallGoodsInfo
	db := global.GVA_DB.Model(&mall.MallGoodsInfo{})
	if keyword != "" {
		db.Where("goods_name like ? or goods_intro like ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if goodsCategoryId >= 0 {
		db.Where("goods_category_id= ?", goodsCategoryId)
	}
	err = db.Count(&total).Error
	switch orderBy {
	case "new":
		db.Order("goods_id desc")
	case "price":
		db.Order("selling_price asc")
	default:
		db.Order("stock_num desc")
	}
	limit := 10
	offset := 10 * (pageNumber - 1)
	err = db.Limit(limit).Offset(offset).Find(&goodsList).Error
	// 返回查询结果
	for _, goods := range goodsList {
		searchGoods := response.GoodsSearchResponse{
			GoodsId:       goods.GoodsId,
			GoodsName:     utils.ReplaceLength(goods.GoodsName, 28),
			GoodsIntro:    utils.ReplaceLength(goods.GoodsIntro, 28),
			GoodsCoverImg: goods.GoodsCoverImg,
			SellingPrice:  goods.SellingPrice,
		}
		searchGoodsList = append(searchGoodsList, searchGoods)
	}
	return
}

```

### 4.8.购物车界面编写

#### 4.8.1.购物车路由接口编写

```golang
package mall

import (
	"Graduation/global"
	"Graduation/model/common/response"
	"Graduation/model/mall/request"
	"Graduation/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MallShopCartApi struct {
}

// 获取购物车列表信息
func (m *MallShopCartApi) CartItemList(c *gin.Context) {
	token := c.GetHeader("token")
	if err, shopCartItem := mallShopCartService.GetShopCartItems(token); err != nil {
		global.GVA_LOG.Error("获取购物车失败", zap.Error(err))
		response.FailWithMessage("获取购物车失败:"+err.Error(), c)
	} else {
		response.OkWithData(shopCartItem, c)
	}
}

// 添加购物车
func (m *MallShopCartApi) AddMallShopCartItem(c *gin.Context) {
	token := c.GetHeader("token")
	var req request.SaveCartItemParam
	_ = c.ShouldBindJSON(&req)
	if err := mallShopCartService.AddMallCartItem(token, req); err != nil {
		global.GVA_LOG.Error("添加购物车失败", zap.Error(err))
		response.FailWithMessage("添加购物车失败:"+err.Error(), c)
	}
	response.OkWithMessage("添加购物车成功", c)
}

// 更新购物车信息
func (m *MallShopCartApi) UpdateMallShopCartItem(c *gin.Context) {
	token := c.GetHeader("token")
	var req request.UpdateCartItemParam
	_ = c.ShouldBindJSON(&req)
	if err := mallShopCartService.UpdateMallCartItem(token, req); err != nil {
		global.GVA_LOG.Error("修改购物车失败", zap.Error(err))
		response.FailWithMessage("修改购物车失败:"+err.Error(), c)
	}
	response.OkWithMessage("修改购物车成功", c)
}

// 删除商品
func (m *MallShopCartApi) DelMallShoppingCartItem(c *gin.Context) {
	token := c.GetHeader("token")
	id, _ := strconv.Atoi(c.Param("newBeeMallShoppingCartItemId"))
	if err := mallShopCartService.DeleteMallCartItem(token, id); err != nil {
		global.GVA_LOG.Error("删除购物车商品失败", zap.Error(err))
		response.FailWithMessage("删除购物车商品失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("删除购物车商品成功", c)
	}
}

// 获取购物信息
func (m *MallShopCartApi) ShopTotal(c *gin.Context) {
	cartItemIdsStr := c.Query("cartItemIds")
	token := c.GetHeader("token")
	cartItemIds := utils.StrToList(cartItemIdsStr)
	if err, cartItemRes := mallShopCartService.GetCartItemsTotal(token, cartItemIds); err != nil {
		global.GVA_LOG.Error("获取购物明细异常：", zap.Error(err))
		response.FailWithMessage("获取购物明细异常:"+err.Error(), c)
	} else {
		response.OkWithData(cartItemRes, c)
	}

}
```

#### 4.8.2.购物车对应增删改查数据库操作

```golang
package mall

import (
	"Graduation/global"
	"Graduation/model/mall"
	"Graduation/model/mall/request"
	"Graduation/model/mall/response"
	"Graduation/utils"
	"errors"

	"github.com/jinzhu/copier"
)

type MallShopCartService struct {
}

// 获取购物车信息列表不分页
func (m *MallShopCartService) GetShopCartItems(token string) (err error, cartItems []response.CartItemResponse) {
	var shopCartItems []mall.MallShopCartItem
	var goodsInfos []mall.MallGoodsInfo
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户"), cartItems
	}
	uuid, _, _ := utils.UndoToken(token)
	global.GVA_DB.Where("u_uid=? and is_deleted = 0", uuid).Find(&shopCartItems)
	var goodsIds []int
	for _, shopcartItem := range shopCartItems {
		goodsIds = append(goodsIds, shopcartItem.GoodsId)
	}
	global.GVA_DB.Where("goods_id in ?", goodsIds).Find(&goodsInfos)
	goodsMap := make(map[int]mall.MallGoodsInfo)
	for _, goodsInfo := range goodsInfos {
		goodsMap[goodsInfo.GoodsId] = goodsInfo
	}
	for _, v := range shopCartItems {
		var cartItem response.CartItemResponse
		copier.Copy(&cartItem, &v)
		if _, ok := goodsMap[v.GoodsId]; ok {
			goodsInfo := goodsMap[v.GoodsId]
			cartItem.GoodsName = goodsInfo.GoodsName
			cartItem.GoodsCoverImg = goodsInfo.GoodsCoverImg
			cartItem.SellingPrice = goodsInfo.SellingPrice
		}
		cartItems = append(cartItems, cartItem)
	}

	return
}

// 添加商品到购物车
func (m *MallShopCartService) AddMallCartItem(token string, req request.SaveCartItemParam) (err error) {
	if req.GoodsCount < 1 {
		return errors.New("商品数量不能小于 1 ！")
	}
	if req.GoodsCount > 5 {
		return errors.New("超出单个商品的最大购买数量！")
	}
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	var shopCartItems []mall.MallShopCartItem
	// 是否已存在商品
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid = ? and goods_id = ? and is_deleted = 0", uuid, req.GoodsId).Find(&shopCartItems).Error
	if err != nil {
		return errors.New("已存在！无需重复添加！")
	}
	err = global.GVA_DB.Where("goods_id = ? ", req.GoodsId).First(&mall.MallGoodsInfo{}).Error
	if err != nil {
		return errors.New("商品为空")
	}
	var total int64
	global.GVA_DB.Where("u_uid = ? and is_deleted = 0", uuid).Count(&total)
	if total > 20 {
		return errors.New("超出购物车最大容量！")
	}
	var shopCartItem mall.MallShopCartItem
	if err = copier.Copy(&shopCartItem, &req); err != nil {
		return err
	}
	shopCartItem.UUid = uuid
	err = global.GVA_DB.Save(&shopCartItem).Error
	return
}

// 更新用户购物车
func (m *MallShopCartService) UpdateMallCartItem(token string, req request.UpdateCartItemParam) (err error) {
	//超出单个商品的最大数量
	if req.GoodsCount > 5 {
		return errors.New("超出单个商品的最大购买数量！")
	}
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	var shopCartItem mall.MallShopCartItem
	if err = global.GVA_DB.Where("cart_item_id=? and is_deleted = 0", req.CartItemId).First(&shopCartItem).Error; err != nil {
		return errors.New("未查询到记录！")
	}
	uuid, _, _ := utils.UndoToken(token)
	if shopCartItem.UUid != uuid {
		return errors.New("未查询到您的信息,禁止该操作！")
	}
	shopCartItem.GoodsCount = req.GoodsCount
	err = global.GVA_DB.Save(&shopCartItem).Error
	return
}

// 删除用户购物车商品
func (m *MallShopCartService) DeleteMallCartItem(token string, id int) (err error) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	var shopCartItem mall.MallShopCartItem
	err = global.GVA_DB.Where("cart_item_id = ? and is_deleted = 0", id).First(&shopCartItem).Error
	if err != nil {
		return
	}
	uuid, _, _ := utils.UndoToken(token)
	if shopCartItem.UUid != uuid {
		return errors.New("未查询到您的信息,禁止该操作！")
	}
	err = global.GVA_DB.Where("cart_item_id = ? and is_deleted = 0", id).UpdateColumns(&mall.MallShopCartItem{IsDeleted: 1}).Error
	return
}

// 获取购物车列表信息总和
func (m *MallShopCartService) GetCartItemsTotal(token string, cartItemIds []int) (err error, cartItemRes []response.CartItemResponse) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户"), cartItemRes
	}
	var shopCartItems []mall.MallShopCartItem
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("cart_item_id in (?) and u_uid = ? and is_deleted = 0", cartItemIds, uuid).Find(&shopCartItems).Error
	if err != nil {
		return
	}
	_, cartItemRes = getMallShopCartItemVOS(shopCartItems)
	//购物车算价
	priceTotal := 0
	for _, cartItem := range cartItemRes {
		priceTotal = priceTotal + cartItem.GoodsCount*cartItem.SellingPrice
	}
	return
}

// 购物车数据转换
func getMallShopCartItemVOS(cartItems []mall.MallShopCartItem) (err error, cartItemsRes []response.CartItemResponse) {
	var goodsIds []int
	for _, cartItem := range cartItems {
		goodsIds = append(goodsIds, cartItem.GoodsId)
	}
	var mallGoods []mall.MallGoodsInfo
	err = global.GVA_DB.Where("goods_id in ?", goodsIds).Find(&mallGoods).Error
	if err != nil {
		return
	}
	mallGoodsMap := make(map[int]mall.MallGoodsInfo)
	for _, goodsInfo := range mallGoods {
		mallGoodsMap[goodsInfo.GoodsId] = goodsInfo
	}
	for _, cartItem := range cartItems {
		var cartItemRes response.CartItemResponse
		copier.Copy(&cartItemRes, &cartItem)
		// 是否包含key
		if _, ok := mallGoodsMap[cartItemRes.GoodsId]; ok {
			mallGoodsTemp := mallGoodsMap[cartItemRes.GoodsId]
			cartItemRes.GoodsCoverImg = mallGoodsTemp.GoodsCoverImg
			goodsName := utils.ReplaceLength(mallGoodsTemp.GoodsName, 28)
			cartItemRes.GoodsName = goodsName
			cartItemRes.SellingPrice = mallGoodsTemp.SellingPrice
			cartItemsRes = append(cartItemsRes, cartItemRes)
		}
	}
	return
}

```

### 4.9. 订单界面

#### 4.9.1.订单路由转接

```golang
package mall

import (
	"Graduation/global"
	"Graduation/model/common/response"
	"Graduation/model/mall/request"
	"Graduation/utils"

	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MallOrderApi struct {
}

// 生成订单
func (m *MallOrderApi) SaveOrder(c *gin.Context) {
	var saveOrderParam request.SaveOrderParam
	_ = c.ShouldBindJSON(&saveOrderParam)
	if err := utils.Verify(saveOrderParam, utils.SaveOrderParamVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	token := c.GetHeader("token")

	priceTotal := 0
	err, itemsForSave := mallShopCartService.GetCartItemsTotal(token, saveOrderParam.CartItemIds)
	if len(itemsForSave) < 1 {
		response.FailWithMessage("无数据:"+err.Error(), c)
	} else {
		//总价
		for _, mallShopCartItemVO := range itemsForSave {
			priceTotal = priceTotal + mallShopCartItemVO.GoodsCount*mallShopCartItemVO.SellingPrice
		}
		if priceTotal < 1 {
			response.FailWithMessage("价格异常", c)
		}
		_, userAddress := mallUserAddressService.GetUserDefaultAddress(token)
		if err, saveOrderResult := mallOrderService.SaveOrder(token, userAddress, itemsForSave); err != nil {
			global.GVA_LOG.Error("生成订单失败", zap.Error(err))
			response.FailWithMessage("生成订单失败:"+err.Error(), c)
		} else {
			response.OkWithData(saveOrderResult, c)
		}
	}
}

// 订单支付
func (m *MallOrderApi) PaySuccess(c *gin.Context) {
	orderNo := c.Query("orderNo")
	payType, _ := strconv.Atoi(c.Query("payType"))
	if err := mallOrderService.PaySuccess(orderNo, payType); err != nil {
		global.GVA_LOG.Error("订单支付失败", zap.Error(err))
		response.FailWithMessage("订单支付失败:"+err.Error(), c)
	}
	response.OkWithMessage("订单支付成功", c)
}

// 完成订单
func (m *MallOrderApi) FinishOrder(c *gin.Context) {
	orderNo := c.Param("orderNo")
	token := c.GetHeader("token")
	if err := mallOrderService.FinishOrder(token, orderNo); err != nil {
		global.GVA_LOG.Error("订单签收失败", zap.Error(err))
		response.FailWithMessage("订单签收失败:"+err.Error(), c)
	}
	response.OkWithMessage("订单签收成功", c)

}

// 取消订单
func (m *MallOrderApi) CancelOrder(c *gin.Context) {
	orderNo := c.Param("orderNo")
	token := c.GetHeader("token")
	if err := mallOrderService.CancelOrder(token, orderNo); err != nil {
		global.GVA_LOG.Error("订单签收失败", zap.Error(err))
		response.FailWithMessage("订单签收失败:"+err.Error(), c)
	}
	response.OkWithMessage("订单签收成功", c)

}

// 订单详情页面
func (m *MallOrderApi) OrderDetailPage(c *gin.Context) {
	orderNo := c.Param("orderNo")
	token := c.GetHeader("token")
	if err, orderDetail := mallOrderService.GetOrderDetailByOrderNo(token, orderNo); err != nil {
		global.GVA_LOG.Error("查询订单详情接口", zap.Error(err))
		response.FailWithMessage("查询订单详情接口:"+err.Error(), c)
	} else {
		response.OkWithData(orderDetail, c)
	}
}

// 订单列表显示
func (m *MallOrderApi) OrderList(c *gin.Context) {
	token := c.GetHeader("token")
	pageNumber, _ := strconv.Atoi(c.Query("pageNumber"))
	status := c.Query("status")
	if pageNumber <= 0 {
		pageNumber = 1
	}
	if err, list, total := mallOrderService.MallOrderListBySearch(token, pageNumber, status); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败"+err.Error(), c)
	} else if len(list) < 1 {
		// 前端项目这里有一个取数逻辑，如果数组为空，数组需要为[] 不能是Null
		response.OkWithDetailed(response.PageResult{
			List:       make([]interface{}, 0),
			TotalCount: total,
			CurrPage:   pageNumber,
			PageSize:   5,
		}, "SUCCESS", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:       list,
			TotalCount: total,
			CurrPage:   pageNumber,
			PageSize:   5,
		}, "SUCCESS", c)
	}

}
```

#### 4.9.2 订单数据库操作

```golang
package mall

import (
	"Graduation/global"
	"Graduation/model/mall"
	"Graduation/model/mall/response"
	"Graduation/model/manage"
	"Graduation/model/manage/request"
	"Graduation/utils"
	"Graduation/utils/enum"
	"errors"

	"github.com/jinzhu/copier"

	"time"
)

type MallOrderService struct {
}

// SaveOrder 保存订单
func (m *MallOrderService) SaveOrder(token string, userAddress mall.MallUserAddress, shopCartItems []response.CartItemResponse) (err error, orderNo string) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户"), orderNo
	}
	uuid, _, _ := utils.UndoToken(token)
	var itemIdList []int
	var goodsIds []int
	for _, cartItem := range shopCartItems {
		itemIdList = append(itemIdList, cartItem.CartItemId)
		goodsIds = append(goodsIds, cartItem.GoodsId)
	}
	var mallGoods []manage.MallGoodsInfo
	global.GVA_DB.Where("goods_id in ? ", goodsIds).Find(&mallGoods)
	//检查是否包含已下架商品
	for _, mallGood := range mallGoods {
		if mallGood.GoodsSellStatus != enum.GOODS_UNDER.Code() {
			return errors.New("已下架，无法生成订单"), orderNo
		}
	}
	mallGoodsMap := make(map[int]manage.MallGoodsInfo)
	for _, mallGood := range mallGoods {
		mallGoodsMap[mallGood.GoodsId] = mallGood
	}
	//判断商品库存
	for _, shopCartItemVO := range shopCartItems {
		//查出的商品中不存在购物车中的这条关联商品数据，直接返回错误提醒
		if _, ok := mallGoodsMap[shopCartItemVO.GoodsId]; !ok {
			return errors.New("购物车数据异常！"), orderNo
		}
		if shopCartItemVO.GoodsCount > mallGoodsMap[shopCartItemVO.GoodsId].StockNum {
			return errors.New("库存不足！"), orderNo
		}
	}
	// 删除购物项
	if len(itemIdList) > 0 && len(goodsIds) > 0 {
		if err = global.GVA_DB.Where("cart_item_id in ?", itemIdList).Updates(mall.MallShopCartItem{IsDeleted: 1}).Error; err == nil {
			var stockNumDTOS []request.StockNumDTO
			copier.Copy(&stockNumDTOS, &shopCartItems)
			for _, stockNumDTO := range stockNumDTOS {
				var goodsInfo manage.MallGoodsInfo
				global.GVA_DB.Where("goods_id =?", stockNumDTO.GoodsId).First(&goodsInfo)
				if err = global.GVA_DB.Where("goods_id =? and stock_num>= ? and goods_sell_status = 0", stockNumDTO.GoodsId, stockNumDTO.GoodsCount).Updates(manage.MallGoodsInfo{StockNum: goodsInfo.StockNum - stockNumDTO.GoodsCount}).Error; err != nil {
					return errors.New("库存不足！"), orderNo
				}
			}
			//生成订单号
			orderNo = utils.GenOrderNo()
			priceTotal := 0
			//保存订单
			var mallOrder manage.MallOrder
			mallOrder.OrderNo = orderNo
			mallOrder.UUid = uuid
			//总价
			for _, mallShopCartItemVO := range shopCartItems {
				priceTotal = priceTotal + mallShopCartItemVO.GoodsCount*mallShopCartItemVO.SellingPrice
			}
			if priceTotal < 1 {
				return errors.New("订单价格异常！"), orderNo
			}
			mallOrder.TotalPrice = priceTotal
			mallOrder.ExtraInfo = ""
			//生成订单项并保存订单项纪录
			if err = global.GVA_DB.Save(&mallOrder).Error; err != nil {
				return errors.New("订单入库失败！"), orderNo
			}
			//生成订单收货地址快照，并保存至数据库
			var mallOrderAddress mall.MallOrderAddress
			copier.Copy(&mallOrderAddress, &userAddress)
			mallOrderAddress.OrderId = mallOrder.OrderId
			//生成所有的订单项快照，并保存至数据库
			var mallOrderItems []manage.MallOrderItem
			for _, mallShoppingCartItemVO := range shopCartItems {
				var mallOrderItem manage.MallOrderItem
				copier.Copy(&mallOrderItem, &mallShoppingCartItemVO)
				mallOrderItem.OrderId = mallOrder.OrderId
				mallOrderItems = append(mallOrderItems, mallOrderItem)
			}
			if err = global.GVA_DB.Save(&mallOrderItems).Error; err != nil {
				return err, orderNo
			}
		}
	}
	return
}

// PaySuccess 支付订单
func (m *MallOrderService) PaySuccess(orderNo string, payType int) (err error) {
	var mallOrder manage.MallOrder
	err = global.GVA_DB.Where("order_no = ? and is_deleted=0 ", orderNo).First(&mallOrder).Error
	if mallOrder != (manage.MallOrder{}) {
		if mallOrder.OrderStatus != 0 {
			return errors.New("订单状态异常！")
		}
		mallOrder.OrderStatus = enum.ORDER_PAID.Code()
		mallOrder.PayType = payType
		mallOrder.PayStatus = 1
		localSH, _ := time.LoadLocation("Asia/Shanghai")
		mallOrder.PayTime = time.Now().In(localSH)
		err = global.GVA_DB.Save(&mallOrder).Error
	}
	return
}

// FinishOrder 完结订单
func (m *MallOrderService) FinishOrder(token string, orderNo string) (err error) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	uuid, _, _ := utils.UndoToken(token)
	var mallOrder manage.MallOrder
	if err = global.GVA_DB.Where("order_no=? and is_deleted = 0", orderNo).First(&mallOrder).Error; err != nil {
		return errors.New("未查询到记录！")
	}
	if mallOrder.UUid != uuid {
		return errors.New("未查询到您的信息,禁止该操作！")
	}
	mallOrder.OrderStatus = enum.ORDER_SUCCESS.Code()
	err = global.GVA_DB.Save(&mallOrder).Error
	return
}

// CancelOrder 关闭订单
func (m *MallOrderService) CancelOrder(token string, orderNo string) (err error) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	uuid, _, _ := utils.UndoToken(token)
	var mallOrder manage.MallOrder
	if err = global.GVA_DB.Where("order_no=? and is_deleted = 0", orderNo).First(&mallOrder).Error; err != nil {
		return errors.New("未查询到记录！")
	}
	if mallOrder.UUid != uuid {
		return errors.New("未查询到您的信息,禁止该操作！")
	}
	if utils.NumsInList(mallOrder.OrderStatus, []int{enum.ORDER_SUCCESS.Code(),
		enum.ORDER_CLOSED_BY_MALLUSER.Code(), enum.ORDER_CLOSED_BY_EXPIRED.Code(), enum.ORDER_CLOSED_BY_JUDGE.Code()}) {
		return errors.New("订单状态异常！")
	}
	mallOrder.OrderStatus = enum.ORDER_CLOSED_BY_MALLUSER.Code()
	err = global.GVA_DB.Save(&mallOrder).Error
	return
}

// GetOrderDetailByOrderNo 获取订单详情
func (m *MallOrderService) GetOrderDetailByOrderNo(token string, orderNo string) (err error, orderDetail response.MallOrderDetailVO) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户"), orderDetail
	}
	uuid, _, _ := utils.UndoToken(token)
	var mallOrder manage.MallOrder
	if err = global.GVA_DB.Where("order_no=? and is_deleted = 0", orderNo).First(&mallOrder).Error; err != nil {
		return errors.New("未查询到记录！"), orderDetail
	}
	if mallOrder.UUid != uuid {
		return errors.New("未查询到您的信息,禁止该操作！"), orderDetail
	}
	var orderItems []manage.MallOrderItem
	err = global.GVA_DB.Where("order_id = ?", mallOrder.OrderId).Find(&orderItems).Error
	if len(orderItems) <= 0 {
		return errors.New("订单项不存在！"), orderDetail
	}

	var mallOrderItemVOS []response.MallOrderItemVO
	copier.Copy(&mallOrderItemVOS, &orderItems)
	copier.Copy(&orderDetail, &mallOrder)
	// 订单状态前端显示为中文
	_, OrderStatusStr := enum.GetMallOrderStatusEnumByStatus(orderDetail.OrderStatus)
	_, payTapStr := enum.GetMallOrderStatusEnumByStatus(orderDetail.PayType)
	orderDetail.OrderStatusString = OrderStatusStr
	orderDetail.PayTypeString = payTapStr
	orderDetail.NewBeeMallOrderItemVOS = mallOrderItemVOS

	return
}

// MallOrderListBySearch 搜索订单
func (m *MallOrderService) MallOrderListBySearch(token string, pageNumber int, status string) (err error, list []response.MallOrderResponse, total int64) {
	// 判断用户是否存在
	if !IsUserExist(token) {
		return errors.New("不存在的用户"), list, total
	}
	uuid, _, _ := utils.UndoToken(token)
	// 根据搜索条件查询
	var mallOrders []manage.MallOrder
	db := global.GVA_DB.Model(&mallOrders)
	if status != "" {
		db.Where("order_status = ?", status)
	}
	err = db.Where("u_uid =? and is_deleted=0 ", uuid).Count(&total).Error
	// 这里前段没有做滚动加载，直接显示全部订单
	// limit := 5
	offset := 5 * (pageNumber - 1)
	err = db.Offset(offset).Order(" order_id desc").Find(&mallOrders).Error
	var orderListVOS []response.MallOrderResponse
	if total > 0 {
		//数据转换 将实体类转成vo
		copier.Copy(&orderListVOS, &mallOrders)
		//设置订单状态中文显示值
		for _, newBeeMallOrderListVO := range orderListVOS {
			_, statusStr := enum.GetMallOrderStatusEnumByStatus(newBeeMallOrderListVO.OrderStatus)
			newBeeMallOrderListVO.OrderStatusString = statusStr
		}
		// 返回订单id
		var orderIds []int
		for _, order := range mallOrders {
			orderIds = append(orderIds, order.OrderId)
		}
		//获取OrderItem
		var orderItems []manage.MallOrderItem
		if len(orderIds) > 0 {
			global.GVA_DB.Where("order_id in ?", orderIds).Find(&orderItems)
			itemByOrderIdMap := make(map[int][]manage.MallOrderItem)
			for _, orderItem := range orderItems {
				itemByOrderIdMap[orderItem.OrderId] = []manage.MallOrderItem{}
			}
			for k, v := range itemByOrderIdMap {
				for _, orderItem := range orderItems {
					if k == orderItem.OrderId {
						v = append(v, orderItem)
					}
					itemByOrderIdMap[k] = v
				}
			}
			//封装每个订单列表对象的订单项数据
			for _, mallOrderListVO := range orderListVOS {
				if _, ok := itemByOrderIdMap[mallOrderListVO.OrderId]; ok {
					orderItemListTemp := itemByOrderIdMap[mallOrderListVO.OrderId]
					var newBeeMallOrderItemVOS []response.MallOrderItemVO
					copier.Copy(&newBeeMallOrderItemVOS, &orderItemListTemp)
					mallOrderListVO.NewBeeMallOrderItemVOS = newBeeMallOrderItemVOS
					_, OrderStatusStr := enum.GetMallOrderStatusEnumByStatus(mallOrderListVO.OrderStatus)
					mallOrderListVO.OrderStatusString = OrderStatusStr
					list = append(list, mallOrderListVO)
				}
			}
		}
	}
	return err, list, total
}
```

## 5. 后台管理系统

### 5.1. 管理员相关基础操作

#### 5.1.1. 管理员的登录和修改信息

```golang
package manage

import (
	"Graduation/global"
	"Graduation/model/common/request"
	"Graduation/model/common/response"
	req "Graduation/model/manage/request"
	"strconv"

	// "Graduation/model/manage/"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ManageAdminUserApi struct {
}

// 管理员用户登录(包括超级管理员)
// AdminLogin 管理员登陆
func (m *ManageAdminUserApi) ManageLogin(c *gin.Context) {
	var manageLoginParams req.ManageLoginParam
	_ = c.ShouldBindJSON(&manageLoginParams)
	if err, msg, token := manageUserService.ManageLogin(manageLoginParams); msg == "Ban" {
		response.FailWithMessage("抱歉,您已被禁用,请联系超级管理员解除!", c)
	} else if err != nil {
		response.FailWithMessage("登陆失败,请检查账号密码是否正确!", c)
	} else {
		response.OkWithData(token, c)
	}
}

// 管理员用户退出(包括超级管理员)
// AdminLogout 管理员登出
func (m *ManageAdminUserApi) ManageLogout(c *gin.Context) {
	token := c.GetHeader("token")
	if err := manageUserService.DeleteManageUserToken(token); err != nil {
		response.FailWithMessage("登出失败", c)
	} else {
		response.OkWithMessage("登出成功", c)
	}
}

// 管理员信息显示
// AdminUserProfile 用id查询AdminUser
func (m *ManageAdminUserApi) ManageUserInfo(c *gin.Context) {
	token := c.GetHeader("token")
	if err, mallAdminUser := manageUserService.GetManageUserInfo(token); err != nil {
		global.GVA_LOG.Error("未查询到管理员信息记录", zap.Error(err))
		response.FailWithMessage("未查询到管理员信息记录", c)
	} else {
		// 扰乱加密,防止泄露
		mallAdminUser.LoginPassword = "******"
		response.OkWithData(mallAdminUser, c)
	}
}

// 修改昵称
func (m *ManageAdminUserApi) UpdateManageUserNickName(c *gin.Context) {
	var reqs req.ManageUpdateNameParam
	_ = c.ShouldBindJSON(&reqs)
	token := c.GetHeader("token")
	if err := manageUserService.UpdateManageUserNickName(token, reqs); err != nil {
		global.GVA_LOG.Error("更新管理员用户昵称失败!", zap.Error(err))
		response.FailWithMessage("更新管理员用户昵称失败", c)
	} else {
		response.OkWithMessage("更新管理员用户昵称成功", c)
	}
}

// 修改密码
func (m *ManageAdminUserApi) UpdateManageUserPassword(c *gin.Context) {
	var reqs req.ManageUpdatePasswordParam
	_ = c.ShouldBindJSON(&reqs)
	userToken := c.GetHeader("token")
	if err := manageUserService.UpdateManagePassWord(userToken, reqs); err != nil {
		global.GVA_LOG.Error("更新密码失败!", zap.Error(err))
		response.FailWithMessage("更新密码失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("更新密码成功", c)
	}
}

// 用户商家列表显示
// UserList 商城注册商家用户列表
func (m *ManageAdminUserApi) UserList(c *gin.Context) {
	var pageInfo req.MallUserSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := manageUserService.GetManageUserInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取管理员用户失败!", zap.Error(err))
		response.FailWithMessage("获取管理员用户失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:       list,
			TotalCount: total,
			CurrPage:   pageInfo.PageNumber,
			PageSize:   pageInfo.PageSize,
		}, "获取管理员用户成功", c)
	}
}

// LockUser 用户禁用[0]与解除禁用[1](0-未锁定 1-已锁定)
func (m *ManageAdminUserApi) LockUser(c *gin.Context) {
	lockStatus, _ := strconv.Atoi(c.Param("lockStatus"))
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := manageUserService.LockUser(IDS, lockStatus); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

```


#### 5.1.2. 管理员增删改查的数据库

```golang
package manage

import (
	"Graduation/global"
	"Graduation/model/common/request"
	"Graduation/model/mall"
	mag "Graduation/model/manage"
	req "Graduation/model/manage/request"
	mallservice "Graduation/service/mall"
	"Graduation/utils"
	"errors"
	"fmt"
	"time"
)

type ManageUserService struct {
}

// 管理员用户登录以及超级管理员登录
func (m *ManageUserService) ManageLogin(req req.ManageLoginParam) (err error, msg string, token string) {
	var mallAdminUser mag.MallAdminUser
	// 管理员就是用户(先查用户表,用户表不存在,则无法登陆,用户表存在则加入管理员账户)
	var user mall.MallUser
	err = global.GVA_DB.Where("login_name=? AND password_md5=?", req.UserName, req.PasswordMd5).First(&user).Error
	if err != nil || user == (mall.MallUser{}) {
		return err, "不存在用户,请注册成为商城用户后才有权限登录", token
	}
	// 查询到用户已存在,
	// 判断是否已经加入到管理员表
	err = global.GVA_DB.Where("login_user_name=? AND login_password=?", req.UserName, req.PasswordMd5).First(&mallAdminUser).Error
	if mallAdminUser == (mag.MallAdminUser{}) {
		// 没加入 则将用户信息加入到管理员信息表super_admin_user
		{
			mallAdminUser.UUid = user.UUid
			mallAdminUser.LoginUserName = user.LoginName
			mallAdminUser.LoginPassword = user.PasswordMd5
			mallAdminUser.NickName = user.NickName
			mallAdminUser.IsSuperAdmin = 0 // 普通管理员

		}
	}
	// 更新管理员表和用户表锁定状态相同
	mallAdminUser.Locked = user.LockedFlag
	// 如果用户被禁,则无法登陆后台管理员系统
	if mallAdminUser.Locked == 1 {
		// 但仍然需要更新用户信息数据
		err = global.GVA_DB.Save(&mallAdminUser).Error
		return err, "Ban", token
	}
	// 创建token
	token, _ = utils.CreateToken(user.UUid)
	strUuid := fmt.Sprintf("%d", user.UUid)
	err = global.GVA_REDIS.Set(global.GVA_CTX, strUuid, token, 3600*time.Second).Err()
	if err != nil {
		global.GVA_LOG.Error("redis存储token失败")
		return err, msg, token
	} else {
		global.GVA_LOG.Info("redis存储token成功")
	}
	err = global.GVA_DB.Save(&mallAdminUser).Error
	return

}

// 管理员以及超级管理员登出,删除管理员登陆token
func (m *ManageUserService) DeleteManageUserToken(token string) (err error) {
	uuid, err, ok := utils.UndoToken(token)
	if err != nil && ok == 0 { // 解码token出现错误
		return err
	}
	uid := fmt.Sprintf("%d", uuid)
	if ok == 1 { // 超时 token已经失效
		global.GVA_REDIS.Del(global.GVA_CTX, uid)
		return nil
	}
	_, err = global.GVA_REDIS.Get(global.GVA_CTX, uid).Result()
	if err != nil {
		global.GVA_LOG.Error("无法从redis得到对应uid的token信息,可能不存在")
		return err
	}
	// 存在则删除
	global.GVA_REDIS.Del(global.GVA_CTX, uid)
	return nil
}

// 检查管理员token是否存在
func (m *ManageUserService) ExistManageToken(token string) (err error, tm int64) {
	uuid, err, ok := utils.UndoToken(token)
	if err != nil && ok == 0 { // 解码token出现错误
		return err, 0
	}
	uid := fmt.Sprintf("%d", uuid)
	_, err = global.GVA_REDIS.Get(global.GVA_CTX, uid).Result()
	if err != nil {
		global.GVA_LOG.Error("token不存在")
		return err, 0
	}
	if ok == 1 {
		global.GVA_LOG.Info("token已超时")
		return err, 1
	}
	return nil, 2
}

// 获取登录管理员和超级管理员的信息
func (m *ManageUserService) GetManageUserInfo(token string) (err error, mallAdminUser mag.MallAdminUser) {
	// 判断用户是否存在
	if !mallservice.IsUserExist(token) {
		return errors.New("不存在的用户"), mallAdminUser
	}
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid = ?", uuid).First(&mallAdminUser).Error
	return err, mallAdminUser
}

// 判断是否是超级管理员
func (m *ManageUserService) IsSuperManageAdmin(token string) (err error, ok bool) {
	// 判断用户是否存在
	if !mallservice.IsUserExist(token) {
		return errors.New("不存在的用户"), false
	}
	var mallAdminUser mag.MallAdminUser
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid = ? And is_super_admin = 1", uuid).First(&mallAdminUser).Error
	if err != nil {
		return err, false
	}
	// 不为空结构体,证明查询到内容
	if mallAdminUser != (mag.MallAdminUser{}) {
		return nil, true
	}
	return nil, false
}

// 更新管理员用户昵称
func (m *ManageUserService) UpdateManageUserNickName(token string, reqs req.ManageUpdateNameParam) (err error) {
	// 判断用户是否存在
	if !mallservice.IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	uuid, _, _ := utils.UndoToken(token)
	// 更新管理员表的昵称
	if err = global.GVA_DB.Where("u_uid = ?", uuid).Updates(&mag.MallAdminUser{
		NickName: reqs.NickName,
	}).Error; err != nil {
		return err
	}

	// 更新用户表里的昵称
	if err = global.GVA_DB.Where("u_uid = ?", uuid).Updates(&mall.MallUser{
		NickName: reqs.NickName,
	}).Error; err != nil {
		return err
	}
	return

}

// 更新管理员用户密码
func (m *ManageUserService) UpdateManagePassWord(token string, reqs req.ManageUpdatePasswordParam) (err error) {
	// 判断用户是否存在
	if !mallservice.IsUserExist(token) {
		return errors.New("不存在的用户")
	}
	var adminUser mag.MallAdminUser
	uuid, _, _ := utils.UndoToken(token)
	err = global.GVA_DB.Where("u_uid =?", uuid).First(&adminUser).Error
	if err != nil {
		return errors.New("不存在的管理员用户")
	}
	if adminUser.LoginPassword != reqs.OriginalPassword {
		return errors.New("原密码不正确")
	}
	if reqs.NewPassword == "" {
		return errors.New("密码不能为空")
	}
	adminUser.LoginPassword = reqs.NewPassword
	// 更新管理员表
	if err = global.GVA_DB.Where("u_uid=?", uuid).Updates(&adminUser).Error; err != nil {
		return
	}
	// 更新用户表
	var userInfo mall.MallUser
	if err = global.GVA_DB.Where("u_uid =?", uuid).First(&userInfo).Error; err != nil {
		return errors.New("管理员用户密码更新失败")
	}
	userInfo.PasswordMd5 = reqs.NewPassword
	err = global.GVA_DB.Save(&userInfo).Error
	return

}

// 查看用户管理员的信息列表
// GetManageUserInfoList 分页获取商城注册用户即管理员列表
func (m *ManageUserService) GetManageUserInfoList(info req.MallUserSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageNumber - 1)
	// 创建db
	db := global.GVA_DB.Model(&mall.MallUser{})
	var mallUsers []mall.MallUser
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&mallUsers).Error
	return err, mallUsers, total
}

// LockUser 超级管理员修改管理员用户状态
func (m *ManageUserService) LockUser(idReq request.IdsReq, lockStatus int) (err error) {
	// 0 正常,1 禁止
	if lockStatus != 0 && lockStatus != 1 {
		return errors.New("操作非法！")
	}
	// 更新 用户表 UpdateColumns locked_flag
	err = global.GVA_DB.Model(&mall.MallUser{}).Where("user_id in ?", idReq.Ids).Update("locked_flag", lockStatus).Error
	return err
}

```

### 5.2. 区块链的使用




