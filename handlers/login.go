package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/spf13/viper"
)

type LoginResponse struct {
	IsAdmin bool   `json:"admin"`
	RollNo  string `json:"rollno"`
}

type LoginRequest struct {
	RollNo   string `json:"rollno"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error decoding request body")
	}

	ok, err := auth.Login(loginRequest.RollNo, loginRequest.Password)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, err.Error())
	}

	if !ok {
		return errors.NewHTTPError(err, http.StatusUnauthorized, "invalid credentials")
	}

	token, err := auth.GenerateToken(loginRequest.RollNo)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
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
		IsAdmin: account.IsAdmin(loginRequest.RollNo),
		RollNo:  loginRequest.RollNo,
	})
	return nil
}

func Logout(w http.ResponseWriter, r *http.Request) error {
	http.SetCookie(w, &http.Cookie{
		Name:     viper.GetString("JWT.COOKIE_NAME"),
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Path:     "/",
	})
	return nil
}
