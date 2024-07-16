package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
}

var Cfg Config

func LoadConfig() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	Cfg = Config{
		DatabaseURL: os.Getenv("MYSQL_URI"),
		ServerPort:  os.Getenv("SERVER_PORT"),
	}
}
