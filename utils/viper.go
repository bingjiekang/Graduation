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
