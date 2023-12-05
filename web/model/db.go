package model

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDb initializes the database connection
func InitDb() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		viper.GetString("db.host"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.dbname"),
		viper.GetString("db.port"),
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
}

// AutoMigrateDb migrates the database
func AutoMigrateDb() {
	err := db.AutoMigrate(&User{}, &Game{}, &Contest{}, &Contestant{}, &Ai{}, &Match{}, &Sdk{})
	if err != nil {
		panic(err)
	}
}
