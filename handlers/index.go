package handlers

import (
	"fmt"
	"net/http"
)

type Response struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error"`
}

func Index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "System is up")
}
