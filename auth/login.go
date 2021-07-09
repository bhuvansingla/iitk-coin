package auth

import (
	"errors"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/util"
)

func Login(u *account.Account) (string, error) {
	if account.ValidateRollNo(u.RollNo) != nil {
		return "", account.ValidateRollNo(u.RollNo)
	}
	if !account.UserExists(u.RollNo) {
		return "", errors.New("account does not exist")
	}
	if !util.CompareHashAndPassword(account.GetStoredPassword(u), u.Password) {
		return "", errors.New("passsword does not match")
	}
	token, err := GenerateToken(u.RollNo)
	if err != nil {
		return "", errors.New("error generating the token")
	}
	return token, nil
}
