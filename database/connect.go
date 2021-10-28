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

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)
    log.Info("Connecting to database", psqlconn)

	psqlconn = "postgres://mkeopogsvhqncb:c1fbd9d5bc2c043e353198280ed0fb1de1402c5b3dbec6ce98927628d3b6bbf1@ec2-52-200-68-5.compute-1.amazonaws.com:5432/d3orbffmu420pg"
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
