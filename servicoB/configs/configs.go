package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port           string `mapstructure:"PORT"`
	WeatherApiUrl  string `mapstructure:"WEATHER_API_URL"`
	ViaCepUrl      string `mapstructure:"VIA_CEP_URL"`
	WeatherApiKey  string `mapstructure:"WEATHER_API_KEY"`
	WeatherBaseURL string `mapstructure:"WEATHER_BASE_URL"`
}

func LoadConfig() (*Config, error) {
	var cfg *Config

	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return cfg, nil
}
