package config

import (
	"github.com/gweebg/probum-users/utils"

	"github.com/spf13/viper"
)

var config *viper.Viper

func Init(env string) {

	config = viper.New()

	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")

	err := config.ReadInConfig()
	utils.Check(err, "")

}

func GetConfig() *viper.Viper {
	return config
}
