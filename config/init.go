package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.yml")
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("unable to locate config file")
	}
}
