package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	config2 "go-nextjs-dashboard/internal/config"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := config2.Cfg.DatabaseURL
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to db", err)
	}

	DB = database
}
