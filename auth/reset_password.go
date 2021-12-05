package auth

import (
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/bhuvansingla/iitk-coin/util"
)

func ResetPassword(rollNo string, newPassword string, otp string) error {

	userExists, err := account.UserExists(rollNo)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if !userExists {
		return errors.NewHTTPError(nil, http.StatusBadRequest, "account doesnot exists")
	}

	if err := account.ValidatePassword(newPassword); err != nil {
		return err
	}

	if err := VerifyOTP(rollNo, otp); err != nil {
		return err
	}

	hashedPwd, err := util.HashAndSalt(newPassword)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	err = account.UpdatePassword(rollNo, hashedPwd)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}
