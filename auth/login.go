package auth

import (
	"errors"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/util"
)

func Login(rollno string, password string) (token string, err error) {
	if err = account.ValidateRollNo(rollno); err != nil {
		return "", err
	}
	if !account.UserExists(rollno) {
		return "", errors.New("account does not exist")
	}
	if !util.CompareHashAndPassword(account.GetStoredPassword(rollno), password) {
		return "", errors.New("passsword does not match")
	}
	if token, err = GenerateToken(rollno); err != nil {
		return "", errors.New("error generating the token")
	}
	return
}
