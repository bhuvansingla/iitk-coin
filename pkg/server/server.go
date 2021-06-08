package server

import (
	"fmt"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/pkg/router"
	_ "github.com/mattn/go-sqlite3"
)

type Auth struct {
	RollNo   string
	Password string
}

func StartServer() {
	router.SetRoutes()
	http.ListenAndServe(":8080", nil)
}

func secretPage(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Secret Page")
}
