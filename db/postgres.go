package db

import (
	"github.com/go-pg/pg"
)

func Connect() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "demo",
		Addr:     "localhost:5432",
	})
	return db
}
