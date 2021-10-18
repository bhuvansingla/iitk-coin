package auth

import (
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/bhuvansingla/iitk-coin/util"
)

func Login(rollno string, password string) (ok bool, err error) {
	if err = account.ValidateRollNo(rollno); err != nil {
		return false, err
	}

	userExists, err := account.UserExists(rollno)

	if err != nil {
		return false, errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if !userExists {
		return false, errors.NewHTTPError(nil, http.StatusBadRequest, "user account does not exist")
	}

	passwordFromRollno, err := account.GetStoredPassword(rollno)

	if err != nil {
		return false, errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if !util.CompareHashAndPassword(passwordFromRollno, password) {
		return false, errors.NewHTTPError(nil, http.StatusBadRequest, "invalid password")
	}
	return true, nil
}
