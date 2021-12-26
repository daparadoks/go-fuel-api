package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

const (
	_configPath = "./resources"
	_configType = "yaml"
)

func NewConfigurationManager() IConfigurationManager {
	env := os.Getenv("PROFILE")
	if env == "" {
		env = "local"
	}

	viper.AddConfigPath(_configPath)
	viper.SetConfigType(_configType)
	applicationConfig := readApplicationConfig(env)
	queConfig := readQueConfig()
	return &ConfigurationManager{applicationConfig: applicationConfig, queConfig: queConfig}
}

type IConfigurationManager interface {
	GetRabbitConfig() RabbitConfig
	GetQuesConfig() QuesConfig
}

type ConfigurationManager struct {
	applicationConfig ApplicationConfig
	queConfig         QuesConfig
}

func (configurationManager *ConfigurationManager) GetRabbitConfig() RabbitConfig {
	return configurationManager.applicationConfig.Rabbit
}

func readApplicationConfig(env string) ApplicationConfig {
	viper.SetConfigName("application")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic("Config load error %s", err.Error())
	}

	var config ApplicationConfig
	c := viper.Sub(env)
	unMarshallErr := c.Unmarshal(&config)
	unMarshallSubErr := c.Unmarshal(&config)
	if unMarshallErr != nil {
		log.Panic("Config load error %s", err.Error())
	}
	if unMarshallSubErr != nil {
		log.Panic("Config load error %s", err.Error())
	}

	return config
}

func readQueConfig() QuesConfig {
	viper.SetConfigName("rabbit-que")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic("Config load error %s", err.Error())
	}

	var config QuesConfig
	c := viper.Sub("que")
	unMarshallErr := c.Unmarshal(&config)
	unMarshallSubErr := c.Unmarshal(&config)

	if unMarshallErr != nil {
		log.Panic("Config load error %s", err.Error())
	}
	if unMarshallSubErr != nil {
		log.Panic("Config load error %s", err.Error())
	}

	return config
}
