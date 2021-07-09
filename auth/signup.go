package auth

import (
	"errors"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/util"
	log "github.com/sirupsen/logrus"
)

func Signup(rollno string, password string, otp string) error {
	if account.UserExists(rollno) {
		return errors.New("account exists already")
	}

	err := account.ValidateRollNo(rollno)
	if err != nil {
		return err
	}

	err = account.ValidatePassword(password)
	if err != nil {
		return err
	}

	ok, err := VerifyOTP(rollno, otp)
	if !ok {
		return errors.New("invalid otp")
	}
	if err != nil {
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

	log.Info("new account created", rollno)
	return nil
}
