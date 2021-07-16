package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
)

type GetCoinBalanceResponse struct {
	Response
	RollNo string `json:"rollno"`
	Coins  int    `json:"coins"`
}

func GetCoinBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	queriedRollno := r.URL.Query().Get("rollno")

	if err := account.ValidateRollNo(queriedRollno); err != nil {
		http.Error(w, "rollno validation failed", http.StatusBadRequest)
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		http.Error(w, "bad cookie", http.StatusBadRequest)
	}

	requestorRole := account.GetAccountRoleByRollno(requestorRollno)

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead || requestorRollno == queriedRollno) {
		http.Error(w, "you are not authorized to read this account balance", http.StatusUnauthorized)
		return
	}

	balance, err := account.GetCoinBalanceByRollno(queriedRollno)

	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(&GetCoinBalanceResponse{
		Coins:  balance,
		RollNo: queriedRollno,
		Response: Response{
			Success: true,
		},
	})
}
