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
