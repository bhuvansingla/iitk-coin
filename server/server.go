package server

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

func StartServer() {
	SetRoutes()
	host := viper.GetString("SERVER.HOST")
	port := viper.GetString("SERVER.PORT")
	address := fmt.Sprintf("%s:%s", host, port)
	http.ListenAndServe(address, nil)
}
