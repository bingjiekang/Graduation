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

// initialization
var snowf, _ = snowflake.GetSnowFlake(0, "", "")

// 生成唯一标识的雪花id
func SnowFlakeUUid() int64 {
	// output ID
	return snowf.Generate()
}


```

#### redis 配置和初始化

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

#### 手机号以及密码验证

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

#### md5 加密

```golang
// 对密码进行加密的md5算法
func Md5(message string) string {
	// 创建一个新的hash对象将字符串转为字节切片
	hash := md5.Sum([]byte(message))
	// 将字节切片转为16进制字符串标识
	return hex.EncodeToString(hash[:])
}
```


#### 登陆、注册、登出代码接口编写

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

#### 登陆、注册、登出代码相关数据库操作编写

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

### 用户地址管理代码编写

#### 增、删、改、查用户地址

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

#### 对应增删改查的数据库操作方法

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

#### 对应修改地址信息的路由建立

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


### 获取商品数据（爬虫）


### 首页信息显示

