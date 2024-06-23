package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	WeatherAPIKey string `mapstructure:"WEATHER_API_KEY"`
}

func LoadConfig(path string) (*Config, error) {
	var conf *Config
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}
