package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type GetNameResponse struct {
	RollNo string   `json:"rollNo"`
	Name   string   `json:"name"`
}

func GetNameByRollNo(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "GET" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	queriedRollNo := r.URL.Query().Get("rollNo")

	if err := account.ValidateRollNo(queriedRollNo); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "invalid rollNo")
	}

	_, err := auth.GetRollNoFromRequest(r)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "invalid cookie")
	}

	name, err := account.GetNameByRollNo(queriedRollNo)

	if err != nil{
		return errors.NewHTTPError(err, http.StatusBadRequest, "could not find username")
	}

	json.NewEncoder(w).Encode(&GetNameResponse{
		Name:  name,
		RollNo: queriedRollNo,
	})
	return nil
}
