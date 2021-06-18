package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDB() error {
	db, err := sql.Open("sqlite3", "iitkcoin.db")
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func CreateUserTable() error {
	_, err := DB.Exec("create table if not exists Account (rollno text, name text, password text, coins int)")
	return err
}
