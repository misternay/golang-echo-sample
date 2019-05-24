package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	User     string
	Host     string
	Password string
	Port     string
	Database string
}

func DB() *DBConfig {
	err := godotenv.Load(".env_dev")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	config := &DBConfig{
		User:     os.Getenv("DB_USERNAME"),
		Host:     os.Getenv("DB_HOST"),
		Password: os.Getenv("DB_PASSWORD"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
	}

	return config
}
