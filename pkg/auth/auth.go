package auth

import (
	"errors"

	jwt "github.com/bhuvansingla/iitk-coin/pkg/jwt"
	"github.com/bhuvansingla/iitk-coin/pkg/user"
	"github.com/bhuvansingla/iitk-coin/pkg/util"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func Login(u *user.User) (string, error) {
	if !user.Exists(u) {
		return "", errors.New("user does not exist")
	}
	if !util.CompareHashAndPassword(user.GetStoredPassword(u), u.Password) {
		return "", errors.New("passsword does not match")
	}
	token, err := jwt.GenerateToken()
	if err != nil {
		return "", errors.New("error generating the token")
	}
	return token, nil
}

func Signup(u *user.User) error {
	if user.Exists(u) {
		return errors.New("user exists already")
	}

	err := user.ValidateRollNo(u)
	if err != nil {
		return err
	}

	err = user.ValidatePassword(u)
	if err != nil {
		return err
	}

	hashedPwd, err := util.HashAndSalt(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPwd

	err = user.Create(u)
	if err != nil {
		return err
	}

	log.Info("new user signed up ", u.RollNo)
	return nil
	// fmt.Fprintf(w, newUser.Password, newUser.RollNo)
}
