package main

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const dbName = "database/db.sqlite3"

func CheckInitDB() *sql.DB {
	if fileExists() {
		db, err := sql.Open("sqlite3", dbName)
		if err != nil {
			panic(err)
		}
		return db

	} else {
		// create file
		file, err := os.Create(dbName)
		if err != nil {
			panic(err)
		}
		file.Close()

		// initialize
		db, err := sql.Open("sqlite3", dbName)
		if err != nil {
			panic(err)
		}

		ddl, err := os.ReadFile("ddl.sql")
		if err != nil {
			panic(err)
		}

		db.Exec(string(ddl))
		return db
	}
}

func fileExists() bool {
	_, err := os.Stat(dbName)
	return err == nil
}
