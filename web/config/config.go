package config

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// InitConfig initializes the configuration of the application
func InitConfig() {
	//这段代码会在web文件夹中生成一个log文件记录输出，如不需要请注释
	f, _ := os.OpenFile("./gin.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	viper.AutomaticEnv()
	// We can use `redis.host` instead of `REDIS_HOST`
	//viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read remaining configs from file
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println(errors.New("config file not found"))
		} else {
			fmt.Println(errors.New("config file was found but another error was produced"))
		}
		return
	}
	viper.WatchConfig()
}
