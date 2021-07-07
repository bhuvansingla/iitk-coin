package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
)

type AddCoinRequest struct {
	Coins  int    `json:"coins"`
	RollNo string `json:"rollno"`
}

type TransferCoinRequest struct {
	Coins      int    `json:"coins"`
	FromRollNo string `json:"fromRollno"`
	ToRollNo   string `json:"toRollno"`
}

type GetCoinBalanceResponse struct {
	Response
	RollNo string `json:"rollno"`
	Coins  int    `json:"coins"`
}

func AddCoins(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		var res Response
		res.Success = false
		res.ErrorMessage = "only POST method allowed"
		json.NewEncoder(w).Encode(res)
		return
	}

	var req AddCoinRequest
	var res Response

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		res.Success = false
		res.ErrorMessage = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestorRole := account.GetAccountRoleByRollno(requestorRollno)
	beneficiaryRole := account.GetAccountRoleByRollno(req.RollNo)

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead || requestorRole == account.CoreTeamMember) {
		res.Success = false
		res.ErrorMessage = "you don't have permission to add coins"
		json.NewEncoder(w).Encode(res)
		return
	}

	if beneficiaryRole == account.GeneralSecretary || beneficiaryRole == account.AssociateHead {
		res.Success = false
		res.ErrorMessage = "beneficiary can't be awarded coins"
		json.NewEncoder(w).Encode(res)
		return
	}

	if beneficiaryRole == account.CoreTeamMember && !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead) {
		res.Success = false
		res.ErrorMessage = "only gensec ah can add coins"
		json.NewEncoder(w).Encode(res)
		return
	}

	err = account.AddCoins(req.RollNo, req.Coins)

	if err != nil {
		res.Success = false
		res.ErrorMessage = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Success = true
	json.NewEncoder(w).Encode(res)
}

func TransferCoins(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		var res Response
		res.Success = false
		res.ErrorMessage = "only POST method allowed"
		json.NewEncoder(w).Encode(res)
		return
	}

	var req TransferCoinRequest
	var res Response

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		res.Success = false
		res.ErrorMessage = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if requestorRollno != req.FromRollNo {
		res.Success = false
		res.ErrorMessage = "send from your own wallet lol"
		json.NewEncoder(w).Encode(res)
		return
	}

	err = account.TransferCoins(req.FromRollNo, req.ToRollNo, req.Coins)

	if err != nil {
		res.Success = false
		res.ErrorMessage = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Success = true
	json.NewEncoder(w).Encode(res)
}

func GetCoinBalance(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		json.NewEncoder(w).Encode(&Response{
			Success:      false,
			ErrorMessage: "only GET method allowed",
		})
		return
	}

	queriedRollno := r.URL.Query().Get("rollno")

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestorRole := account.GetAccountRoleByRollno(requestorRollno)

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead || requestorRollno == queriedRollno) {
		json.NewEncoder(w).Encode(&Response{
			Success:      false,
			ErrorMessage: "you are not authorized to check this wallet ballance",
		})
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