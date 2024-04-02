package config

import (
	"github.com/spf13/viper"
)

var C = &Config{}

// Config 配置
type Config struct {
	UserName      string `toml:"UserName" yaml:"username"`
	Password      string `toml:"Password" yaml:"password"`
	ConcurrentNum int    `toml:"ConcurrentNum" yaml:"concurrent_mum"` // 并发数
}

func ParseConfig(filePath string) {
	viper.SetConfigType("toml")
	viper.AddConfigPath(filePath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			panic(err)
		} else {
			// Config file was found but another error was produced
			panic(err)
		}
	}
	if err := viper.Unmarshal(C); err != nil {
		panic(err)
	}
}
