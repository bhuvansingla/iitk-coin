package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/pkg/auth"
	"github.com/bhuvansingla/iitk-coin/pkg/user"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func Login(w http.ResponseWriter, req *http.Request) {
	var u user.User
	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	token, err := auth.Login(&u)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, token)
}

func Signup(w http.ResponseWriter, req *http.Request) {

	var u user.User
	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		log.Error(err.Error())
		fmt.Fprint(w, err.Error())
		return
	}
	err = auth.Signup(&u)
	if err != nil {
		log.Error(err.Error())
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, "signed up successfully")
}
