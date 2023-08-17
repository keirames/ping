package config

import (
	"main/logger"

	"github.com/spf13/viper"
)

type config struct {
	Env          string `mapstructure:"ENV"`
	DBSource     string `mapstructure:"DB_SOURCE"`
	DBDriverName string `mapstructure:"DB_DRIVER_NAME"`
	DatabaseURL  string `mapstructure:"DATABASE_URL"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	Port         string `mapstructure:"PORT"`
}

var C config

func Load() error {
	viper.SetDefault("Env", "dev")
	viper.SetDefault("PORT", "8080")
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
