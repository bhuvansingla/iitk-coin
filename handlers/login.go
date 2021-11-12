package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/spf13/viper"
)

type LoginResponse struct {
	IsAdmin bool   `json:"admin"`
	RollNo  string `json:"rollNo"`
}

type LoginRequest struct {
	RollNo   string `json:"rollNo"`
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
		return err
	}

	if !ok {
		return errors.NewHTTPError(err, http.StatusUnauthorized, "invalid credentials")
	}

	return setCookiesAndRespond(loginRequest.RollNo, w)
}

func setCookiesAndRespond(rollNo string, w http.ResponseWriter) error {

	accessToken, err := auth.GenerateAccessToken(rollNo)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	refreshToken, err := auth.GenerateRefreshToken(rollNo)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	setCookies(w, accessToken, refreshToken)

	isAdmin, err := account.IsAdmin(rollNo)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	err = json.NewEncoder(w).Encode(&LoginResponse{
		IsAdmin: isAdmin,
		RollNo:  rollNo,
	})
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}

func setCookies(w http.ResponseWriter, accessToken string, refreshToken string) {
	cookie := &http.Cookie{
		Name:     viper.GetString("JWT.ACCESS_TOKEN.NAME"),
		Value:    accessToken,
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(w, cookie)

	cookie = &http.Cookie{
		Name:     viper.GetString("JWT.REFRESH_TOKEN.NAME"),
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/auth",
	}

	http.SetCookie(w, cookie)
}

func Logout(w http.ResponseWriter, r *http.Request) error {

	cookie, err := r.Cookie(viper.GetString("JWT.REFRESH_TOKEN.NAME"))
	if err != nil {
		setCookies(w, "", "")
		return nil
	}

	rollNo, err := auth.GetRollNoFromTokenCookie(cookie)
	if err != nil {
		setCookies(w, "", "")
		return nil
	}

	setCookies(w, "", "")

	err = account.DeleteToken(rollNo)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}
