package config

import (
	"github.com/spf13/viper"
	"log"
)

type DatabaseConfiguration struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     uint16
}

type MinioConfiguration struct {
	Host      string
	Port      string
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type Configuration struct {
	Database DatabaseConfiguration
	Minio    MinioConfiguration
}

func ProvideConfig() Configuration {
	var config Configuration

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalln(err)
	}

	return config
}
