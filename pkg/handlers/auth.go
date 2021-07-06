package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bhuvansingla/iitk-coin/pkg/account"
	"github.com/bhuvansingla/iitk-coin/pkg/auth"
	"github.com/bhuvansingla/iitk-coin/pkg/cors"
	"github.com/bhuvansingla/iitk-coin/pkg/jwt"
	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

type SignupResponse struct {
	Response
}

type LoginResponse struct {
	Response
	IsAdmin bool   `json:"admin"`
	RollNo  string `json:"rollno"`
}

func Login(w http.ResponseWriter, req *http.Request) {

	cors.SetPolicy(&w, req)
	w.Header().Set("Content-Type", "application/json")

	if req.Method != "POST" {
		var res Response
		res.Success = false
		res.ErrorMessage = "only POST method allowed"
		json.NewEncoder(w).Encode(res)
		return
	}

	var res LoginResponse
	var u account.Account

	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		res.Success = false
		res.ErrorMessage = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	token, err := auth.Login(&u)
	if err != nil {
		res.Success = false
		res.ErrorMessage = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	log.Info(cookie)
	res.Success = true
	res.IsAdmin = account.IsAdmin(u.RollNo)
	res.RollNo = u.RollNo
	json.NewEncoder(w).Encode(res)
}

func Signup(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if req.Method != "POST" {
		var res Response
		res.Success = false
		res.ErrorMessage = "only POST method allowed"
		json.NewEncoder(w).Encode(res)
		return
	}

	var res SignupResponse
	var u account.Account

	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		res.Success = false
		res.ErrorMessage = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	err = auth.Signup(&u)
	if err != nil {
		res.Success = false
		res.ErrorMessage = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Success = true
	json.NewEncoder(w).Encode(res)
}

func CheckUserIsLoggedIn(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if req.Method != "POST" {
		var res Response
		res.Success = false
		res.ErrorMessage = "only POST method allowed"
		json.NewEncoder(w).Encode(res)
		return
	}

	requestorRollno, err := jwt.GetRollnoFromRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var res LoginResponse
	res.Success = true
	res.IsAdmin = account.IsAdmin(requestorRollno)
	res.RollNo = requestorRollno
	json.NewEncoder(w).Encode(res)
}
