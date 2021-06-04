package main

import (
	"database/sql"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	RollNo string
	Name   string
}

var privateKey = []byte("keyy")

func addUser(db *sql.DB, user *User) {
	if userExists(db, user) {
		fmt.Println("User", user.RollNo, "Exists Already")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error begining transaction")
	}
	stmt, err := db.Prepare("INSERT INTO User (rollno,name) VALUES (?,?)")
	if err != nil {
		fmt.Println("Error preparing statement")
	}
	_, err = stmt.Exec(user.RollNo, user.Name)
	if err != nil {
		fmt.Println("Error executing statement")
	}
	tx.Commit()
	fmt.Println("User", user.RollNo, "Added Succesfully")
}

func userExists(db *sql.DB, user *User) bool {
	row := db.QueryRow("SELECT rollno FROM User WHERE rollno=?", user.RollNo)
	scannedRow := ""
	row.Scan(&scannedRow)
	return scannedRow != ""
}

func index(w http.ResponseWriter, req *http.Request) {

	token, err := generateToken()

	if err != nil {
	}
	fmt.Println(token)
	fmt.Fprintf(w, "Yayyyyy!")
	fmt.Println("Works!")
}

func secretPage(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Secret Page")

}

var db *sql.DB

func generateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = false
	tokenString, err := token.SignedString(privateKey)

	if err != nil {
		fmt.Println("Error here", err.Error())
		return "", err
	}

	return tokenString, nil
}

func isAuth(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
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

func main() {

	// db, errOpenDB := sql.Open("sqlite3", "coin.db")

	// if errOpenDB != nil {
	// 	fmt.Println("Error opening DB")
	// }

	// db.Exec("create table if not exists User (rollno text, name text)")

	// addUser(db, &User{RollNo: "180199", Name: "Bhuvan Singla"})
	// addUser(db, &User{RollNo: "180199", Name: "Bhuvan Singla"})

	http.HandleFunc("/", index)
	http.Handle("/access", isAuth(secretPage))
	http.ListenAndServe(":8080", nil)
	fmt.Println("Listening on 8080")
}
