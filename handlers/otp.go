package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type OtpRequest struct {
	RollNo string `json:"rollNo"`
}

func GenerateOtp(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	var otpRequest OtpRequest
	if err := json.NewDecoder(r.Body).Decode(&otpRequest); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error decoding request body")
	}

	if err := auth.GenerateOtp(otpRequest.RollNo); err != nil {
		return err
	}

	return nil
}
