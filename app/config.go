package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var cfgFile string
var config *Config

type Config struct {
}

func initConfig() error {
	viper.SetConfigFile(cfgFile)
	viper.SetEnvPrefix("xpto")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		logrus.WithField("message", err.Error()).Errorln("read config file")
		return err
	}
	fmt.Println("using config file:", viper.ConfigFileUsed())
	err := viper.Unmarshal(&config)
	if err != nil {
		logrus.WithField("message", err.Error()).Errorln("decode config file")
		return err
	}
	return nil
}
