package server

import (
	"net/http"
)

func StartServer() {
	SetRoutes()
	http.ListenAndServe(":8080", nil)
}
