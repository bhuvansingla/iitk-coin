package auth

import (
	"errors"

	"github.com/bhuvansingla/iitk-coin/pkg/account"
	jwt "github.com/bhuvansingla/iitk-coin/pkg/jwt"
	"github.com/bhuvansingla/iitk-coin/pkg/util"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func Login(u *account.Account) (string, error) {
	if account.ValidateRollNo(u) != nil {
		return "", account.ValidateRollNo(u)
	}
	if !account.UserExists(u.RollNo) {
		return "", errors.New("account does not exist")
	}
	if !util.CompareHashAndPassword(account.GetStoredPassword(u), u.Password) {
		return "", errors.New("passsword does not match")
	}
	token, err := jwt.GenerateToken(u.RollNo)
	if err != nil {
		return "", errors.New("error generating the token")
	}
	return token, nil
}

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
