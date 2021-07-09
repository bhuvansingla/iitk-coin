package account

import (
	"errors"

	"github.com/bhuvansingla/iitk-coin/db"
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

func Create(rollno string, hashedPasssword string, name string) error {

	role := NormalUser
	stmt, err := db.DB.Prepare("INSERT INTO Account (rollno,name,password,coins,role) VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(rollno, name, hashedPasssword, 0, role)
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

func ValidateRollNo(rollno string) error {
	if rollno == "" {
		return errors.New("empty roll no")
	}
	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("empty password")
	}
	if len(password) < 8 {
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
