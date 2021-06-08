package jwt

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
)

var privateKey = []byte("SHHHHH!! SECRET HAI!")

func GenerateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["rollno"] = false
	tokenString, err := token.SignedString(privateKey)

	if err != nil {
		fmt.Println("Error here", err.Error())
		return "", err
	}

	return tokenString, nil
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Invalid Signing Method")
				}

				checkAdmin := token.Claims.(jwt.MapClaims)["admin"]
				if checkAdmin == true {
					fmt.Printf("Admin")
				} else {
					return nil, fmt.Errorf("Not Admin")
				}

				return privateKey, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "No Auth Token")
		}
	})
}
