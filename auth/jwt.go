package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

var privateKey = []byte(viper.GetString("JWT.PRIVATE_KEY"))

type Claims struct {
	Rollno string `json:"rollno"`
	jwt.RegisteredClaims
}

func GenerateToken(rollno string) (string, error) {

	expirationTime := time.Now().Add(time.Duration(viper.GetInt("JWT.EXPIRATION_TIME_IN_MIN")) * time.Minute)

	claims := &Claims{
		Rollno: rollno,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(privateKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(viper.GetString("JWT.COOKIE_NAME"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad token"))
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return privateKey, nil
		})

		//check time

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad token"))
			return
		}

		if token.Valid {
			endpoint(w, r)
			return
		}
	}
}

func GetRollnoFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie(viper.GetString("JWT.COOKIE_NAME"))
	if err != nil {
		return "", err
	}
	return GetRollnoFromTokenCookie(cookie)
}

func GetRollnoFromTokenCookie(cookie *http.Cookie) (string, error) {
	token := cookie.Value
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	if err != nil {
		return "", err
	}
	return claims.Rollno, nil
}
