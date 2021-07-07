package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	log "github.com/sirupsen/logrus"
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

	auth.SetCorsPolicy(&w, req)
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

	requestorRollno, err := auth.GetRollnoFromRequest(req)
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

func GenerateOtp(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if req.Method != "POST" {
		json.NewEncoder(w).Encode(&Response{
			Success:      false,
			ErrorMessage: "only POST method allowed",
		})
		return
	}

	requestorRollno, err := auth.GetRollnoFromRequest(req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = auth.GenerateOtp(requestorRollno)

	if err != nil {
		var res Response
		res.Success = false
		res.ErrorMessage = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	json.NewEncoder(w).Encode(&Response{
		Success: true,
	})
}
