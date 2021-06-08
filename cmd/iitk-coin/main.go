package main

import (
	"fmt"

	"github.com/bhuvansingla/iitk-coin/pkg/db"
	"github.com/bhuvansingla/iitk-coin/pkg/server"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := db.ConnectDB()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = db.CreateUserTable()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	server.StartServer()
}
