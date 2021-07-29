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
	if !account.UserExists(rollno) {
		return false, errors.NewHTTPError(nil, http.StatusBadRequest, "account does not exist")
	}
	if !util.CompareHashAndPassword(account.GetStoredPassword(rollno), password) {
		return false, errors.NewHTTPError(nil, http.StatusBadRequest, "invalid password")
	}
	return true, nil
}
