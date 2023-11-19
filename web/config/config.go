package config

import (
	"strings"

	"github.com/spf13/viper"
)

// InitConfig initializes the configuration of the application
func InitConfig() {
	viper.AutomaticEnv()
	// We can use `redis.host` instead of `REDIS_HOST`
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read remaining configs from file
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	viper.WatchConfig()
}
