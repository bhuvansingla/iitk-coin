package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type GetCoinBalanceResponse struct {
	RollNo string `json:"rollno"`
	Coins  int    `json:"coins"`
}

func GetCoinBalance(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "GET" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	queriedRollno := r.URL.Query().Get("rollno")

	if err := account.ValidateRollNo(queriedRollno); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "Invalid rollno")
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "Invalid cookie")
	}

	requestorRole, err := account.GetAccountRoleByRollno(requestorRollno)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error when getting account role")
	}

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead || requestorRollno == queriedRollno) {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "You are not authorized to read this account balance")
	}

	userExists, err := account.UserExists(queriedRollno)

	if err != nil {
		return errors.NewHTTPError(nil, http.StatusBadRequest, "error when checking if account exists")
	}

	if !userExists {
		return errors.NewHTTPError(err, http.StatusBadRequest, "account does not exist")
	}

	balance, err := account.GetCoinBalanceByRollno(queriedRollno)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	json.NewEncoder(w).Encode(&GetCoinBalanceResponse{
		Coins:  balance,
		RollNo: queriedRollno,
	})
	return nil
}
