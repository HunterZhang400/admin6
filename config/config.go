package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Cfg *Config

func Setup() {
	viper.AutomaticEnv()
	env := os.Getenv("NODE_ENV")
	if env == "dev" || env == "" {
		viper.AddConfigPath("./config")
		viper.SetConfigName("config.develop")
	} else {
		viper.AddConfigPath("./config")
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Print(err.Error())
		return
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		fmt.Print(err.Error())
		return
	}
	//特殊处理
	//cli, err := aes.NewAESGCM(aes.DefaultCode)
	//if err != nil {
	//	panic(err)
	//}
	//res, err := cli.DecryptUseStandardBase64(Cfg.Mysql.Password)
	//if err != nil {
	//	panic(err)
	//}
	//Cfg.Mysql.Password = string(res)
	Cfg.Server.Env = env

}

type Config struct {
	Server      *Server
	Mysql       *Database
	WoodeBoxAPI *WoodeBoxAPIConfig `mapstructure:"woodebox"`
}
type Server struct {
	Env       string
	JWTSecret string `mapstructure:"jwt-secret"`
}
type Database struct {
	Engine   string
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	MaxIdle  int `mapstructure:"max-idle"`
	MaxOpen  int `mapstructure:"max-open"`
}

type WoodeBoxAPIConfig struct {
	BaseURL    string `mapstructure:"base-url"`
	Timeout    int    `mapstructure:"timeout"`
	MaxRetries int    `mapstructure:"max-retries"`
	Token      string `mapstructure:"token"`
}
