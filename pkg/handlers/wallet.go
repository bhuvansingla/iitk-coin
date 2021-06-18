package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/pkg/account"

	_ "github.com/mattn/go-sqlite3"
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

type GetCoinBalanceRequest struct {
	RollNo string `json:"rollno"`
}

type GetCoinBalanceResponse struct {
	Response
	RollNo string `json:"rollno"`
	Coins  int    `json:"coins"`
}

func AddCoins(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

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
	w.WriteHeader(http.StatusOK)

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
	w.WriteHeader(http.StatusOK)

	if r.Method != "GET" {
		var res Response
		res.Success = false
		res.ErrorMessage = "only GET method allowed"
		json.NewEncoder(w).Encode(res)
		return
	}

	var req GetCoinBalanceRequest
	var res GetCoinBalanceResponse

	res.RollNo = req.RollNo

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		res.Success = false
		res.ErrorMessage = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	balance, err := account.GetCoinBalanceByRollno(req.RollNo)
	if err != nil {
		res.Success = false
		res.ErrorMessage = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Coins = balance
	res.Success = true
	json.NewEncoder(w).Encode(res)
}
