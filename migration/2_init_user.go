package main

import (
	"fmt"

	"echo-sample/handler"
	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		hashPwd, _ := handler.HashPassword("123456")
		var statement = fmt.Sprintf(`INSERT INTO users (username, fullname, password) VALUES ('admin', 'John Due', '%s')`, hashPwd)
		_, err := db.Exec(statement)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("Truncating users...")
		_, err := db.Exec(`TRUNCATE users`)
		return err
	})
}
