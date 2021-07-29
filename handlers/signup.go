package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type SignupRequest struct {
	Rollno   string `json:"rollno"`
	Password string `json:"password"`
	Otp      string `json:"otp"`
}

func Signup(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	var signupRequest SignupRequest

	if err := json.NewDecoder(r.Body).Decode(&signupRequest); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error decoding request body")
	}

	err := auth.Signup(signupRequest.Rollno, signupRequest.Password, signupRequest.Otp)

	if err != nil {
		return err
	}

	return nil
}
