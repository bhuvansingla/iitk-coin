package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

var privateKey = []byte("SHHHHH!! SECRET HAI!")

type Claims struct {
	Rollno string `json:"rollno"`
	jwt.StandardClaims
}

type Response struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error"`
}

func GenerateToken(rollno string) (string, error) {

	expirationTime := time.Now().Add(10 * time.Minute)

	claims := &Claims{
		Rollno: rollno,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
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
		cookie, err := r.Cookie("token")
		if err != nil {
			log.Error(err)
			json.NewEncoder(w).Encode(&Response{
				Success:      false,
				ErrorMessage: "couldn't find cookie token",
			})
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
			log.Error(err)
			json.NewEncoder(w).Encode(&Response{
				Success:      false,
				ErrorMessage: "couldn't parse token",
			})
			return
		}

		if token.Valid {
			endpoint(w, r)
			return
		}
	}
}

func GetRollnoFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
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
