package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type RewardRequest struct {
	Coins	int    `json:"coins"`
	RollNo	string `json:"rollNo"`
	Remarks	string `json:"remarks"`
}

type RewardResponse struct {
	TxnId string `json:"id"`
}

func RewardCoins(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	var rewardRequest RewardRequest

	if err := json.NewDecoder(r.Body).Decode(&rewardRequest); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error decoding request body")
	}

	requestorRollNo, err := auth.GetRollNoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "invalid cookie")
	}

	requestorRole, err := account.GetAccountRoleByRollNo(requestorRollNo)
	if err != nil {
		return err
	}

	beneficiaryRole, err := account.GetAccountRoleByRollNo(rewardRequest.RollNo)
	if err != nil {
		return err
	}

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead || requestorRole == account.CoreTeamMember) {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "you don't have permission to add coins to the requested acccount")
	}

	if beneficiaryRole == account.GeneralSecretary || beneficiaryRole == account.AssociateHead {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "not possible to add coins to the requested acccount")
	}

	if beneficiaryRole == account.CoreTeamMember && !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead) {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "only gensec/ah can add coins to this account")
	}

	id, err := account.AddCoins(rewardRequest.RollNo, rewardRequest.Coins, rewardRequest.Remarks)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(&RewardResponse{TxnId: id})
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}
