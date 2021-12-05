package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type ResetPasswordRequest struct {
	RollNo      string `json:"rollNo"`
	NewPassword string `json:"newPassword"`
	Otp         string `json:"otp"`
}

func ResetPassword(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	var resetPasswordRequest ResetPasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&resetPasswordRequest); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error decoding request body")
	}

	err := auth.ResetPassword(resetPasswordRequest.RollNo, resetPasswordRequest.NewPassword, resetPasswordRequest.Otp)

	if err != nil {
		return err
	}

	return nil
}
