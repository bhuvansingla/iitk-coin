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
	user := viper.GetString("DATABASE.USER")
	password := viper.GetString("DATABASE.PASS")
	dataSource := fmt.Sprintf("%s?_auth&_auth_user=%s&_auth_pass=%s", name, user, password)

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
