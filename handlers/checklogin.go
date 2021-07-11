package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
)

type CheckLoginResponse struct {
	Response
	IsAdmin bool   `json:"admin"`
	RollNo  string `json:"rollno"`
}

func CheckLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		http.Error(w, "bad cookie", http.StatusBadRequest)
	}

	isAdmin := account.IsAdmin(requestorRollno)
	json.NewEncoder(w).Encode(&CheckLoginResponse{
		RollNo: requestorRollno,
		Response: Response{
			Success: true,
		},
		IsAdmin: isAdmin,
	})
}
