package main

import (
	"fmt"

	"github.com/bhuvansingla/iitk-coin/db"
	"github.com/bhuvansingla/iitk-coin/server"
)

func main() {
	err := db.ConnectDB()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = db.CreateTables()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	server.StartServer()
}
