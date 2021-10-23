package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var DB *sql.DB

func Connect() (err error) {
	var (
		host     = viper.GetString("DATABASE.HOST")
		port     = viper.GetString("DATABASE.PORT")
		user     = viper.GetString("DATABASE.USER")
		password = viper.GetString("DATABASE.PASSWORD")
		dbname   = viper.GetString("DATABASE.NAME")
	)

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    
	db, err := sql.Open("postgres", psqlconn)
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
