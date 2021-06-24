package handlers

import (
	"fmt"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/pkg/jwt"
)

type Response struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error"`
}

func Index(w http.ResponseWriter, req *http.Request) {
	rollno, err := jwt.GetRollnoFromRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	fmt.Fprint(w, "Hi! "+rollno)
}
