package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func Start() (err error) {

	setRoutes()

	host := viper.GetString("SERVER.HOST")

	port := os.Getenv("PORT")
	if port == "" {
		port = viper.GetString("SERVER.PORT")
	}
	address := fmt.Sprintf("%s:%s", host, port)

	err = http.ListenAndServe(address, nil)

	return err
}
