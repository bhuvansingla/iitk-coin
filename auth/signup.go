package auth

import (
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/bhuvansingla/iitk-coin/util"
	log "github.com/sirupsen/logrus"
)

func Signup(rollNo string, name string, password string, otp string) error {
	
	userExists, err := account.UserExists(rollNo)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if userExists {
		return errors.NewHTTPError(nil, http.StatusBadRequest, "account exists already")
	}

	if err := account.ValidateRollNo(rollNo); err != nil {
		return err
	}

	if err := account.ValidatePassword(password); err != nil {
		return err
	}

	if err := VerifyOTP(rollNo, otp); err != nil {
		return err
	}

	hashedPwd, err := util.HashAndSalt(password)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	err = account.Create(rollNo, hashedPwd, name)
	if err != nil {
		return err
	}

	log.Info("A new account was created with roll no ", rollNo)

	return nil
}
