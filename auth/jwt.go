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
	RollNo string `json:"rollNo"`
	jwt.RegisteredClaims
}

func GenerateToken(rollNo string) (string, error) {

	expirationTime := time.Now().Add(time.Duration(viper.GetInt("JWT.EXPIRATION_TIME_IN_MIN")) * time.Minute)

	claims := &Claims{
		RollNo: rollNo,
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

		if token.Valid {
			endpoint(w, r)
			return
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("bad token"))
				return
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("token expired"))
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
				return
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}
		
	}
}

func GetRollNoFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie(viper.GetString("JWT.COOKIE_NAME"))
	if err != nil {
		return "", err
	}
	return GetRollNoFromTokenCookie(cookie)
}

func GetRollNoFromTokenCookie(cookie *http.Cookie) (string, error) {
	token := cookie.Value
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	if err != nil {
		return "", err
	}
	return claims.RollNo, nil
}
