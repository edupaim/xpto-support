package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig ServerConfig
}

func (c *Config) WithPort(port int) {
	c.ServerConfig.Port = port
}

type ServerConfig struct {
	Port int
}

func InitConfig(filePath string) (*Config, error) {
	var currentConfig *Config
	err := loadConfigFromFile(filePath, currentConfig)
	if err != nil {
		return nil, err
	}
	viper.SetEnvPrefix("xpto")
	viper.AutomaticEnv()
	return currentConfig, nil
}

func loadConfigFromFile(filePath string, currentConfig *Config) error {
	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Errorln("read config file")
		return err
	}
	fmt.Println("using config file:", viper.ConfigFileUsed())
	err := viper.Unmarshal(currentConfig)
	if err != nil {
		logrus.WithError(err).Errorln("decode config file")
		return err
	}
	return nil
}
