package account

import (
	"database/sql"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/database"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type Role int64

const (
	NormalUser       Role = 0
	GeneralSecretary Role = 1
	AssociateHead    Role = 2
	CoreTeamMember   Role = 3
)

func Create(rollNo string, hashedPasssword string, name string) error {

	role := NormalUser
	stmt, err := database.DB.Prepare("INSERT INTO ACCOUNT (rollNo, name, password, coins, role) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	_, err = stmt.Exec(rollNo, name, hashedPasssword, 0, role)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func UserExists(rollNo string) (bool, error) {
	row := database.DB.QueryRow("SELECT rollNo FROM ACCOUNT WHERE rollNo=$1", rollNo)
	scannedRow := ""
	err := row.Scan(&scannedRow)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return scannedRow != "", nil
}

func GetAccountRoleByRollNo(rollNo string) (Role, error) {
	row := database.DB.QueryRow("SELECT role FROM ACCOUNT WHERE rollNo=$1", rollNo)
	var role Role
	err := row.Scan(&role)

	if err == sql.ErrNoRows {
		return role, errors.NewHTTPError(err, http.StatusBadRequest, "account doesnot exist")
	}
	if err != nil {
		return role, errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return role, nil
}

func GetNameByRollNo(rollNo string) (string, error) {
	row := database.DB.QueryRow("SELECT name FROM ACCOUNT WHERE rollNo=$1", rollNo)
	name := ""
	err := row.Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func GetStoredPassword(rollNo string) (string, error) {
	row := database.DB.QueryRow("SELECT password FROM ACCOUNT WHERE rollNo=$1", rollNo)
	scannedRow := ""
	err := row.Scan(&scannedRow)
	if err != nil {
		return "", err
	}
	return scannedRow, nil
}

func ValidateRollNo(rollNo string) error {
	if rollNo == "" {
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

func IsAdmin(rollNo string) (bool, error) {
	role, err := GetAccountRoleByRollNo(rollNo)
	if err != nil {
		return false, err
	}
	if role == GeneralSecretary || role == AssociateHead || role == CoreTeamMember {
		return true, nil
	}
	return false, nil
}
