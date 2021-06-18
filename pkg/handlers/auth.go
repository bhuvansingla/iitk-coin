package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/pkg/account"
	"github.com/bhuvansingla/iitk-coin/pkg/auth"

	_ "github.com/mattn/go-sqlite3"
)

type SignupResponse struct {
	Response
}

type LoginResponse struct {
	Response
	JwtToken string `json:"token"`
}

func Login(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

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

	res.Success = true
	res.JwtToken = token
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
