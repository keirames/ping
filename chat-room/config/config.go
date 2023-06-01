package config

import (
	"chatroom/logger"

	"github.com/spf13/viper"
)

type config struct {
	DBSource  string `mapstructure:"DB_SOURCE"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
}

var C config

func Load() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		logger.L.Error().Err(err).Msg("Load config env file failed!")
		return err
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		logger.L.Error().Err(err).Msg("Load config env file failed!")
		return err
	}

	return nil
}
