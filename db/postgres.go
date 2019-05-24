package db

import (
	"fmt"
	"os"

	"github.com/go-pg/pg"
	"github.com/joho/godotenv"
)

func Connect() *pg.DB {
	var err error

	if os.Getenv("ENV") == "production" {
		err = godotenv.Load(".env")
	} else {
		err = godotenv.Load(".env_dev")
	}

	if err != nil {
		fmt.Println("Error loading .env or .env_dev file")
	}

	db := pg.Connect(&pg.Options{
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
		Addr:     string(os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")),
	})

	return db
}
