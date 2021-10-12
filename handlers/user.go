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

	queriedRollno := r.URL.Query().Get("rollno")

	if err := account.ValidateRollNo(queriedRollno); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "invalid rollno")
	}

	_, err := auth.GetRollnoFromRequest(r)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "invalid cookie")
	}

	name, err := account.GetNameByRollNo(queriedRollno)

	if err != nil{
		return errors.NewHTTPError(err, http.StatusBadRequest, "could not find username")
	}

	json.NewEncoder(w).Encode(&GetNameResponse{
		Name:  name,
		RollNo: queriedRollno,
	})
	return nil
}
