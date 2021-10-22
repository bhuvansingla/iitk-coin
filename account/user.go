package account

import (
	"net/http"

	"github.com/bhuvansingla/iitk-coin/database"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type Role int

const (
	NormalUser       Role = 0
	GeneralSecretary Role = 1
	AssociateHead    Role = 2
	CoreTeamMember   Role = 3
)

func Create(rollno string, hashedPasssword string, name string) error {

	role := NormalUser
	stmt, err := database.DB.Prepare("INSERT INTO ACCOUNT (rollno,name,password,coins,role) VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(rollno, name, hashedPasssword, 0, role)
	if err != nil {
		return err
	}
	return nil
}

func UserExists(rollno string) (bool, error) {
	row := database.DB.QueryRow("SELECT rollno FROM ACCOUNT WHERE rollno=?", rollno)
	scannedRow := ""
	err := row.Scan(&scannedRow)
	if err != nil {
		return false, err
	}
	return scannedRow != "", nil
}

func GetAccountRoleByRollno(rollno string) (Role, error) {
	row := database.DB.QueryRow("SELECT role FROM ACCOUNT WHERE rollno=?", rollno)
	var role Role
	err := row.Scan(&role)
	if err != nil {
		return NormalUser, err
	}
	return role, nil
}

func GetNameByRollNo(rollno string) (string, error) {
	row := database.DB.QueryRow("SELECT name FROM ACCOUNT WHERE rollno=?", rollno)
	name := ""
	err := row.Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func GetStoredPassword(rollno string) (string, error) {
	row := database.DB.QueryRow("SELECT password FROM ACCOUNT WHERE rollno=?", rollno)
	scannedRow := ""
	err := row.Scan(&scannedRow)
	if err != nil {
		return "", err
	}
	return scannedRow, nil
}

func ValidateRollNo(rollno string) error {
	if rollno == "" {
		return errors.NewHTTPError(nil, http.StatusBadRequest, "roll no is empty")
	}
	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.NewHTTPError(nil, http.StatusBadRequest, "password is empty")
	}
	if len(password) < 8 {
		return errors.NewHTTPError(nil, http.StatusBadRequest, "password is less than 8 characters.")
	}
	return nil
}

func IsAdmin(rollno string) (bool, error) {
	role, err := GetAccountRoleByRollno(rollno)
	if err != nil {
		return false, errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	if role == GeneralSecretary || role == AssociateHead || role == CoreTeamMember {
		return true, nil
	}
	return false, nil
}
