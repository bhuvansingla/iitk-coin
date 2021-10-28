package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type CheckLoginResponse struct {
	IsAdmin bool   `json:"admin"`
	RollNo  string `json:"rollNo"`
}

func CheckLogin(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	requestorRollNo, err := auth.GetRollNoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "invalid cookie")
	}

	isAdmin, err := account.IsAdmin(requestorRollNo)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	json.NewEncoder(w).Encode(&CheckLoginResponse{
		RollNo:  requestorRollNo,
		IsAdmin: isAdmin,
	})

	return nil
}
