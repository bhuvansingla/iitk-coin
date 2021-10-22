package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type WalletHistoryResponse struct {
	History []interface{}	`json:"history"`
	RollNo	string			`json:"rollno"`
}

func WalletHistory(w http.ResponseWriter, r *http.Request) error {

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

	requestorRole := account.GetAccountRoleByRollno(requestorRollno)

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead || requestorRollno == queriedRollno) {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "You are not authorized to read this account history")
	}

	if !account.UserExists(queriedRollno) {
		return errors.NewHTTPError(nil, http.StatusBadRequest, "user account does not exist")
	}

	history, err := account.GetWalletHistoryByRollNo(queriedRollno)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	json.NewEncoder(w).Encode(&WalletHistoryResponse{
		History:	history,
		RollNo:		queriedRollno,
	})
	return nil
}
