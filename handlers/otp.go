package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/auth"
)

type OtpRequest struct {
	Rollno string `json:"rollno"`
}

func GenerateOtp(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var otpRequest OtpRequest

	if err := json.NewDecoder(r.Body).Decode(&otpRequest); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	_, err := auth.GenerateOtp(otpRequest.Rollno)

	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(&Response{
		Success: true,
	})
}
