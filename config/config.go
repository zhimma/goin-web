package config

import (
	"bytes"
	_ "embed"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/zhimma/goin-web/pkg/env"
	"github.com/zhimma/goin-web/pkg/file"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Mysql struct {
		Read struct {
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Database string `yaml:"database"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `json:"read" yaml:"read"`
		Write struct {
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Database string `yaml:"database"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `json:"write" yaml:"write"`
		Base struct {
			MaxOpenConn     int           `yaml:"maxOpenConn"`
			MaxIdleConn     int           `yaml:"maxIdleConn"`
			ConnMaxLifeTime time.Duration `yaml:"connMaxLifeTime"`
		} `json:"dbBase" yaml:"dbBase"`
	} `json:"mysql" yaml:"mysql"`
	Redis struct {
		Addr         string `json:"addr" yaml:"addr"`
		Pass         string `json:"pass" yaml:"pass"`
		Db           int    `json:"db" yaml:"db"`
		MaxRetries   int    `json:"maxRetries" yaml:"maxRetries"`
		PoolSize     int    `json:"poolSize" yaml:"poolSize"`
		MinIdleConns int    `json:"minIdleConns" yaml:"minIdleConns"`
	} `json:"redis" yaml:"redis"`

	Language struct {
		Local string `json:"local" yaml:"local"`
	} `json:"language" yaml:"language"`
}

var config = new(Config)

var (
	//go:embed local.yaml
	localConfig []byte

	//go:embed dev.yaml
	devConfig []byte

	//go:embed uat.yaml
	uatConfig []byte

	//go:embed prod.yaml
	prodConfig []byte
)

func init() {
	var r io.Reader
	switch env.NowEnv().Value() {
	case "local":
		r = bytes.NewReader(localConfig)
	case "dev":
		r = bytes.NewReader(devConfig)
	case "uat":
		r = bytes.NewReader(uatConfig)
	case "prod":
		r = bytes.NewReader(prodConfig)
	default:
		r = bytes.NewReader(devConfig)
	}
	// 设置env文件类型
	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(r); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}

	viper.SetConfigName(env.NowEnv().Value())

	viper.AddConfigPath("./config")

	configFile := "./config/" + env.NowEnv().Value() + ".yaml"
	// 判断文件是否存在
	_, ok := file.IsExists(configFile)
	if !ok {
		if err := os.MkdirAll(filepath.Dir(configFile), 0766); err != nil {
			panic(err)
		}
		f, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := viper.WriteConfig(); err != nil {
			panic(err)
		}

	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})
}

func Get() Config {
	return *config
}
