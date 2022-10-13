package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	config Config
	once   sync.Once
)

func InitializeConfig() *viper.Viper {
	// 设置配置文件路径
	configFile := "config.yaml"
	// 生产环境可以通过设置环境变量来改变配置文件路径
	if configEnv := os.Getenv("VIPER_CONFIG"); configEnv != "" {
		configFile = configEnv
	}

	// 初始化 viper
	v := viper.New()

	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %s", err))
	}

	// 监听配置文件
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)
		// 重载配置
		if err := v.Unmarshal(&config); err != nil {
			fmt.Println(err)
		}
	})

	// 配置赋值
	if err := v.Unmarshal(&config); err != nil {
		fmt.Println(err)
	}

	return v
}

func GetConfig() Config {
	once.Do(func() {
		InitializeConfig()
	})

	return config
}
