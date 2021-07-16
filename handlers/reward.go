package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
)

type RewardRequest struct {
	Coins  int    `json:"coins"`
	RollNo string `json:"rollno"`
}

func RewardCoins(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var rewardRequest RewardRequest

	if err := json.NewDecoder(r.Body).Decode(&rewardRequest); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		http.Error(w, "bad cookie", http.StatusBadRequest)
	}

	requestorRole := account.GetAccountRoleByRollno(requestorRollno)
	beneficiaryRole := account.GetAccountRoleByRollno(rewardRequest.RollNo)

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead || requestorRole == account.CoreTeamMember) {
		http.Error(w, "you don't have permission to add coins to the requested acccount", http.StatusUnauthorized)
		return
	}

	if beneficiaryRole == account.GeneralSecretary || beneficiaryRole == account.AssociateHead {
		http.Error(w, "not possible to add coins to the requested acccount", http.StatusUnauthorized)
		return
	}

	if beneficiaryRole == account.CoreTeamMember && !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead) {
		http.Error(w, "only gensec/ah can add coins to this account", http.StatusUnauthorized)
		return
	}

	err = account.AddCoins(rewardRequest.RollNo, rewardRequest.Coins)

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
