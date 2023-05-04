package configx

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Init
/*
	初始化日志
*/
func Init(configFile string, c any) {
	// 指定配置文件
	viper.SetConfigFile(configFile)
	// 读取配置文件内容
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	// 配置文件内容映射至struct
	if err := viper.Unmarshal(c); err != nil {
		panic(err)
	}
	// 监控配置文件内容
	viper.WatchConfig()
	// 当配置文件内容改变时, 自动重新映射配置文件内容至struct
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(c); err != nil {
			panic(err)
		}
	})
}
