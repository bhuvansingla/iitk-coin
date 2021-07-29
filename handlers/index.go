package handlers

import (
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "System is up")
}
