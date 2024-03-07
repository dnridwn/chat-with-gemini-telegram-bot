package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectToDB() (*sql.DB, error) {
	var (
		host   = os.Getenv("DB_HOST")
		port   = os.Getenv("DB_PORT")
		user   = os.Getenv("DB_USER")
		pass   = os.Getenv("DB_PASS")
		dbName = os.Getenv("DB_NAME")
	)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbName))
	if err != nil {
		return db, err
	}

	return db, nil
}
