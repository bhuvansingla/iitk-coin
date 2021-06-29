package account

import (
	"errors"

	"github.com/bhuvansingla/iitk-coin/pkg/db"
	_ "github.com/mattn/go-sqlite3"
)

type Role int

const (
	NormalUser       Role = 0
	GeneralSecretary Role = 1
	AssociateHead    Role = 2
	CoreTeamMember   Role = 3
)

type Account struct {
	RollNo   string
	Name     string
	Password string
	Coins    int
	Role     Role
}

func Create(account *Account) error {

	account.Role = NormalUser
	stmt, err := db.DB.Prepare("INSERT INTO Account (rollno,name,password,coins,role) VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(account.RollNo, account.Name, account.Password, account.Coins, account.Role)
	if err != nil {
		return err
	}
	return nil
}

func UserExists(rollno string) bool {
	row := db.DB.QueryRow("SELECT rollno FROM Account WHERE rollno=?", rollno)
	scannedRow := ""
	row.Scan(&scannedRow)
	return scannedRow != ""
}

func GetAccountRoleByRollno(rollno string) Role {
	row := db.DB.QueryRow("SELECT role FROM Account WHERE rollno=?", rollno)
	var role Role
	row.Scan(&role) // handle error
	return role
}

func GetStoredPassword(account *Account) string {
	row := db.DB.QueryRow("SELECT password FROM Account WHERE rollno=?", account.RollNo)
	scannedRow := ""
	row.Scan(&scannedRow)
	return scannedRow
}

func ValidateRollNo(account *Account) error {
	if account.RollNo == "" {
		return errors.New("empty roll no")
	}
	return nil
}

func ValidatePassword(account *Account) error {
	if account.Password == "" {
		return errors.New("empty password")
	}
	if len(account.Password) < 8 {
		return errors.New("password too small")
	}
	return nil
}

func IsAdmin(rollno string) bool {
	role := GetAccountRoleByRollno(rollno)
	if role == GeneralSecretary || role == AssociateHead || role == CoreTeamMember {
		return true
	}
	return false
}
