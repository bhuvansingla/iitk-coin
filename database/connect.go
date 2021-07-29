package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var DB *sql.DB

func Connect() (err error) {

	name := viper.GetString("DATABASE.NAME")
	dataSource := fmt.Sprintf("%s", name)

	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		log.Error(err.Error())
		return
	}

	DB = db

	err = createTables()
	if err != nil {
		return
	}

	return
}
