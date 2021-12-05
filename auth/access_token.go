package auth

import (
	"net/http"
	"time"

	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/spf13/viper"
)

func GenerateAccessToken(rollNo string) (string, error) {

	expirationTime := time.Now().Add(time.Duration(viper.GetInt("JWT.ACCESS_TOKEN.EXPIRATION_TIME_IN_MIN")) * time.Minute)

	return generateToken(rollNo, expirationTime)
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(viper.GetString("JWT.ACCESS_TOKEN.NAME"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad token"))
			return
		}

		err = isTokenValid(cookie)

		if err == nil {
			endpoint(w, r)
			return
		}

		errors.WriteResponse(err, w)
	}
}
