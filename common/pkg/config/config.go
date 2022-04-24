package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
)

var conf = new(Config)

func Init(path string) (*Config, error) {
	//默认使用dev.yaml配置文件
	if path == "" {
		viper.AddConfigPath("./conf")
		viper.SetConfigName("dev")
		viper.SetConfigType("yaml")
	} else {
		//parse指定yaml
		viper.SetConfigFile(fmt.Sprintf("conf/%s", path))
	}

	//加载配置
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	//解析
	if err := viper.Unmarshal(conf); err != nil {
		return nil, errors.Wrapf(err, "faild to load config")
	}

	//开启配置热加载
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("config file changed:", in.Name)
	})
	return conf, nil
}

func GetConfig() *Config {
	return conf
}
