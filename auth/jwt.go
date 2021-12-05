package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

var privateKey = []byte(viper.GetString("JWT.PRIVATE_KEY"))

type Claims struct {
	RollNo string `json:"rollNo"`
	jwt.RegisteredClaims
}

func GetRollNoFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie(viper.GetString("JWT.ACCESS_TOKEN.NAME"))
	if err != nil {
		return "", err
	}
	return GetRollNoFromTokenCookie(cookie)
}

func GetRollNoFromTokenCookie(cookie *http.Cookie) (string, error) {
	token := cookie.Value
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, keyFunc)
	if err != nil {
		return "", err
	}
	return claims.RollNo, nil
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("invalid signing method")
	}
	return privateKey, nil
}

func generateToken(rollNo string, expirationTime time.Time) (string, error) {
	
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

func isTokenValid(cookie *http.Cookie) error {

	token, err := jwt.Parse(cookie.Value, keyFunc)

	if token.Valid {
		return nil
	} 
	
	jwtError, ok := err.(*jwt.ValidationError)
	
	if ok {
		if jwtError.Errors&jwt.ValidationErrorMalformed != 0 {
			return errors.NewHTTPError(err, http.StatusBadRequest, "validation malformed")
		} else if jwtError.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return errors.NewHTTPError(err, http.StatusUnauthorized, "token expired")
		} else {
			return errors.NewHTTPError(nil, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	} else {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}
