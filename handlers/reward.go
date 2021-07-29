package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type RewardRequest struct {
	Coins  int    `json:"coins"`
	RollNo string `json:"rollno"`
}

func RewardCoins(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	var rewardRequest RewardRequest

	if err := json.NewDecoder(r.Body).Decode(&rewardRequest); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error decoding request body")
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "Invalid cookie")
	}

	requestorRole := account.GetAccountRoleByRollno(requestorRollno)
	beneficiaryRole := account.GetAccountRoleByRollno(rewardRequest.RollNo)

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead || requestorRole == account.CoreTeamMember) {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "you don't have permission to add coins to the requested acccount")
	}

	if beneficiaryRole == account.GeneralSecretary || beneficiaryRole == account.AssociateHead {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "not possible to add coins to the requested acccount")
	}

	if beneficiaryRole == account.CoreTeamMember && !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead) {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "only gensec/ah can add coins to this account")
	}

	if err = account.AddCoins(rewardRequest.RollNo, rewardRequest.Coins); err != nil {
		return err
	}

	return nil
}
