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
	RollNo  string `json:"rollno"`
}

func CheckLogin(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "Invalid cookie")
	}

	isAdmin, err := account.IsAdmin(requestorRollno)

	if err != nil {
		return err
	}

	json.NewEncoder(w).Encode(&CheckLoginResponse{
		RollNo:  requestorRollno,
		IsAdmin: isAdmin,
	})

	return nil
}
