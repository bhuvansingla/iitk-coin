package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/auth"
)

type SignupRequest struct {
	Rollno   string `json:"rollno"`
	Password string `json:"password"`
	Otp      string `json:"otp"`
}

func Signup(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var signupRequest SignupRequest

	err := json.NewDecoder(r.Body).Decode(&signupRequest)
	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}

	err = auth.Signup(signupRequest.Rollno, signupRequest.Password, signupRequest.Otp)

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
