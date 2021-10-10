package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type GetNameResponse struct {
	RollNo string `json:"rollno"`
	Name  string    `json:"name"`
}

func GetNameByRollno(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "GET" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "Invalid cookie")
	}

	name := account.GetNameByRollNo(requestorRollno)

	json.NewEncoder(w).Encode(&GetNameResponse{
		Name:  name,
		RollNo: requestorRollno,
	})
	return nil
}
