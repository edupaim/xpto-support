package app

import (
	"edupaim/xpto-support/app/domain"
	"edupaim/xpto-support/app/services"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig     ServerConfig
	ArangoConfig     services.ArangoConfig
	LegacyXptoConfig LegacyXptoConfig
}

func (c *Config) WithPort(port int) {
	c.ServerConfig.Port = port
}

type ServerConfig struct {
	Port int
}

type LegacyXptoConfig struct {
	Address string
}

func InitConfig(filePath string, debug bool) (*Config, error) {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	var currentConfig Config
	err := loadConfigFromFile(filePath, &currentConfig)
	if err != nil {
		return nil, err
	}
	viper.SetEnvPrefix("xpto")
	viper.AutomaticEnv()
	passphrase := viper.GetString("passphrase")
	if passphrase != "" {
		err = domain.SetCryptPassphrase(passphrase)
		if err != nil {
			return nil, err
		}
	}
	return &currentConfig, nil
}

func loadConfigFromFile(filePath string, currentConfig *Config) error {
	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Errorln("read config file")
		return err
	}
	logrus.Debugln("using config file:", viper.ConfigFileUsed())
	err := viper.Unmarshal(&currentConfig)
	if err != nil {
		logrus.WithError(err).Errorln("decode config file")
		return err
	}
	logrus.WithField("config", currentConfig).Debugln("current config")
	return nil
}
