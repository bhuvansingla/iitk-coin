package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/spf13/viper"
)

type LoginResponse struct {
	Response
	IsAdmin bool   `json:"admin"`
	RollNo  string `json:"rollno"`
}

type LoginRequest struct {
	RollNo   string `json:"rollno"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	token, err := auth.Login(loginRequest.RollNo, loginRequest.Password)
	if err != nil {
		json.NewEncoder(w).Encode(&LoginResponse{
			Response: Response{
				Success:      false,
				ErrorMessage: err.Error(),
			},
		})
		return
	}

	cookie := &http.Cookie{
		Name:     viper.GetString("JWT.COOKIE_NAME"),
		Value:    token,
		Expires:  time.Now().Add(time.Duration(viper.GetInt("JWT.EXPIRATION_TIME_IN_MIN")) * time.Minute),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	json.NewEncoder(w).Encode(&LoginResponse{
		Response: Response{
			Success: true,
		},
		IsAdmin: account.IsAdmin(loginRequest.RollNo),
		RollNo:  loginRequest.RollNo,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     viper.GetString("JWT.COOKIE_NAME"),
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Path:     "/",
	})
}
