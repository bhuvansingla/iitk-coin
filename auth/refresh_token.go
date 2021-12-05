package auth

import (
	"net/http"
	"time"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/spf13/viper"
)

func GenerateRefreshToken(rollNo string) (string, error) {

	expirationTime := time.Now().Add(time.Duration(viper.GetInt("JWT.REFRESH_TOKEN.EXPIRATION_TIME_IN_MIN")) * time.Minute)

	refreshToken, err := generateToken(rollNo, expirationTime)
	if err != nil {
		return "", err
	}

	err = account.UpdateRefreshToken(refreshToken, rollNo)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func CheckRefreshTokenValidity(r *http.Request) (string, error) {

	cookie, err := r.Cookie(viper.GetString("JWT.ACCESS_TOKEN.NAME"))
	if err != nil {
		return "", errors.NewHTTPError(err, http.StatusBadRequest, "bad access token")
	}

	err = isTokenValid(cookie)

	if err == nil {
		rollNo, err := GetRollNoFromTokenCookie(cookie)
		if err != nil {
			return "", errors.NewHTTPError(err, http.StatusBadRequest, "bad access token")
		}
		return rollNo, nil
	}

	clientError, ok := err.(errors.ClientError)
	if !ok {
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if status, _ := clientError.ResponseHeaders(); status!=http.StatusUnauthorized {
		return "", err
	}

	cookie, err = r.Cookie(viper.GetString("JWT.REFRESH_TOKEN.NAME"))
	if err != nil {
		return "", errors.NewHTTPError(err, http.StatusBadRequest, "bad refresh token")
	}

	rollNo, err := GetRollNoFromTokenCookie(cookie)
	if err != nil {
		return "", errors.NewHTTPError(err, http.StatusBadRequest, "bad refresh token")
	}

	refreshToken, err := account.GetRefreshToken(rollNo)
	if err != nil {
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if refreshToken != cookie.Value {
		return "", errors.NewHTTPError(err, http.StatusBadRequest, "bad refresh token")
	}

	err = isTokenValid(cookie)

	return rollNo, err
}
