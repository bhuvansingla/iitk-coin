package main

import (
	"database/sql"
	"fmt"

	// "net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	rollno string
	name   string
}

func addUser(db *sql.DB, user *User) {
	if userExists(db, user) {
		fmt.Println("User", user.rollno, "Exists Already")
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
	_, err = stmt.Exec(user.rollno, user.name)
	if err != nil {
		fmt.Println("Error executing statement")
	}
	tx.Commit()
	fmt.Println("User", user.rollno, "Added Succesfully")
}

func userExists(db *sql.DB, user *User) bool {
	row := db.QueryRow("SELECT rollno FROM User WHERE rollno=?", user.rollno)
	scannedRow := ""
	row.Scan(&scannedRow)
	return scannedRow != ""
}

// func foo(w http.ResponseWriter, req *http.Request) {
// 	fmt.Println("Works!")
// }

func main() {

	db, errOpenDB := sql.Open("sqlite3", "coin.db")

	if errOpenDB != nil {
		fmt.Println("Error opening DB")
	}

	db.Exec("create table if not exists User (rollno text, name text)")

	addUser(db, &User{rollno:"180199", name:"Bhuvan Singla"})
	addUser(db, &User{rollno:"180199", name:"Bhuvan Singla"})

	// http.HandleFunc("/", foo)
	// http.ListenAndServe(":8080", nil)
}
