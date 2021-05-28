package main

import (
	"database/sql"
	"fmt"

	// "net/http"

	_ "github.com/mattn/go-sqlite3"
)

func addUser(db *sql.DB, rollno string, name string) {
	if userExists(db, rollno) {
		fmt.Println("User", rollno, "Exists Already")
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
	_, err = stmt.Exec(rollno, name)
	if err != nil {
		fmt.Println("Error executing statement")
	}
	tx.Commit()
	fmt.Println("User", rollno, "Added Succesfully")
}

func userExists(db *sql.DB, rollno string) bool {
	row := db.QueryRow("SELECT rollno FROM User WHERE rollno=?", rollno)
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

	addUser(db, "180199", "Bhuvan Singla")
	addUser(db, "180199", "Bhuvan Singla")

	// http.HandleFunc("/", foo)
	// http.ListenAndServe(":8080", nil)
}
