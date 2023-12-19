package config

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"hiper-backend/user"
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
}

// InitConfigAfter initializes the configuration of the application after other modules are initialized
func InitConfigAfter() {
	viper.WatchConfig()
	viper.OnConfigChange(handleConfigChange)
	applyConfig()
}

// handleConfigChange handles changes in the configuration
func handleConfigChange(e fsnotify.Event) {
	log.Println("Config file changed:", e.Name)
	applyConfig()
}

func applyConfig() {
	if viper.IsSet("superadmin.password") {
		user.UpsertSuperAdmin()
	}
}
