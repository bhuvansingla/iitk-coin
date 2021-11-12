package handlers

import (
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "GET" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	rollNo, err := auth.CheckRefreshTokenValidity(r)

	if err != nil {
		return err
	}

	return setCookiesAndRespond(rollNo, w)
}

func InvalidateRefreshTokens(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	requestorRollNo, err := auth.GetRollNoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "invalid cookie")
	}

	requestorRole, err := account.GetAccountRoleByRollNo(requestorRollNo)
	if err != nil {
		return err
	}

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead || requestorRole == account.CoreTeamMember) {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "you don't have permission to invalidate refresh tokens")
	}

	err = account.InvalidateAllTokens()
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}
