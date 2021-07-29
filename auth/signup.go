package auth

import (
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/bhuvansingla/iitk-coin/util"
	log "github.com/sirupsen/logrus"
)

func Signup(rollno string, password string, otp string) error {
	if account.UserExists(rollno) {
		return errors.NewHTTPError(nil, http.StatusBadRequest, "account exists already")
	}

	if err := account.ValidateRollNo(rollno); err != nil {
		return err
	}

	if err := account.ValidatePassword(password); err != nil {
		return err
	}

	if err := VerifyOTP(rollno, otp); err != nil {
		return err
	}

	hashedPwd, err := util.HashAndSalt(password)
	if err != nil {
		return err
	}

	err = account.Create(rollno, hashedPwd, "User Name")
	if err != nil {
		return err
	}

	log.Info("A new account was created with roll no ", rollno)

	return nil
}
