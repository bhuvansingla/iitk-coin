package main

import (
	_ "github.com/bhuvansingla/iitk-coin/config"
	"github.com/bhuvansingla/iitk-coin/database"
	_ "github.com/bhuvansingla/iitk-coin/logger"
	"github.com/bhuvansingla/iitk-coin/server"
)

func main() {

	err := database.Connect()
	if err != nil {
		return
	}

	err = server.Start()
	if err != nil {
		return
	}

}
