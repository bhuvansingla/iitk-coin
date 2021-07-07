package auth

import (
	"errors"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/util"
	log "github.com/sirupsen/logrus"
)

func Signup(u *account.Account) error {
	if account.UserExists(u.RollNo) {
		return errors.New("account exists already")
	}

	err := account.ValidateRollNo(u)
	if err != nil {
		return err
	}

	err = account.ValidatePassword(u)
	if err != nil {
		return err
	}

	hashedPwd, err := util.HashAndSalt(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPwd

	err = account.Create(u)
	if err != nil {
		return err
	}

	log.Info("new account created", u.RollNo)
	return nil
}
