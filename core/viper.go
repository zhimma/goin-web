package core

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	globalInstance "github.com/zhimma/goin-web/global"
)

func Viper() *viper.Viper {
	v := viper.New()
	v.SetConfigFile("config.yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件失败：%s\n", err))
	}
	if err := v.Unmarshal(&globalInstance.BaseConfig); err != nil {
		fmt.Println(err)
	}
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件发生变动 ", in.Name)
		if err := v.Unmarshal(&globalInstance.BaseConfig); err != nil {
			fmt.Println(err)
		}
	})
	return v
}
