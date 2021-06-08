package user

import (
	"errors"

	"github.com/bhuvansingla/iitk-coin/pkg/db"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	RollNo   string
	Name     string
	Password string
	Coins    int
}

func Create(user *User) error {

	stmt, err := db.DB.Prepare("INSERT INTO User (rollno,name,password,coins) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.RollNo, user.Name, user.Password, user.Coins)
	if err != nil {
		return err
	}
	return nil
}

func Exists(user *User) bool {
	row := db.DB.QueryRow("SELECT rollno FROM User WHERE rollno=?", user.RollNo)
	scannedRow := ""
	row.Scan(&scannedRow)
	return scannedRow != ""
}

func GetStoredPassword(user *User) string {
	row := db.DB.QueryRow("SELECT password FROM User WHERE rollno=?", user.RollNo)
	scannedRow := ""
	row.Scan(&scannedRow)
	return scannedRow
}

func ValidateRollNo(u *User) error {
	if u.RollNo == "" {
		return errors.New("empty roll no")
	}
	return nil
}

func ValidatePassword(u *User) error {
	if u.Password == "" {
		return errors.New("empty password")
	}
	if len(u.Password) < 8 {
		return errors.New("password too small")
	}
	return nil
}
