package server

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

func Start() (err error) {

	setRoutes()

	host := viper.GetString("SERVER.HOST")
	port := viper.GetString("SERVER.PORT")
	address := fmt.Sprintf("%s:%s", host, port)

	err = http.ListenAndServe(address, nil)

	return err
}
