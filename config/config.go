package config

import (
	"errors"

	"github.com/spf13/viper"
)

func LoadConfig(file string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigFile(file)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}
