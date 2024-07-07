package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	WEATHER_API_KEY             string `mapstructure:"WEATHER_API_KEY"`
	OTEL_SERVICE_NAME           string `mapstructure:"OTEL_SERVICE_NAME"`
	OTEL_EXPORTER_OTLP_ENDPOINT string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
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
