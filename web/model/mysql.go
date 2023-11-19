package model

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDb initializes the database connection
func InitDb() {
	host := viper.GetString("db.host")
	port := viper.GetString("db.port")
	name := viper.GetString("db.name")
	user := viper.GetString("db.user")
	password := viper.GetString("db.password")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
	// Connect to database
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

// AutoMigrateDb migrates the database
func AutoMigrateDb() {
	// TODO: implement
}
