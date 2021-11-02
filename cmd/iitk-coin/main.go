package main

import (
	_ "github.com/bhuvansingla/iitk-coin/config"
	"github.com/bhuvansingla/iitk-coin/database"
	_ "github.com/bhuvansingla/iitk-coin/logger"
	"github.com/bhuvansingla/iitk-coin/mail"
	"github.com/bhuvansingla/iitk-coin/server"
	log "github.com/sirupsen/logrus"
)

func main() {

	err := database.Connect()
	if err != nil {
		log.Error("Error connecting to database: %s", err)
		return
	}
	
	err =  mail.Test()
	if err != nil {
		log.Error("Error sending mail: %s", err)
		return
	}

	err = server.Start()
	if err != nil {
		log.Error("Error starting server: %s", err)
		return
	}

}
