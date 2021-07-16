package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
)

type NewRedeemRequest struct {
	NumCoins       int    `json:"numCoins"`
	ReceiverRollno string `json:"receiverRollno"`
	Item           string `json:"item"`
	Otp            string `json:"otp"`
}

type UpdateRedeemRequest struct {
	RedeemId int `json:"redeemId"`
}

type RedeemListResponse struct {
	Response
	RedeemList []account.RedeemRequest `json:"redeemList`
}

func NewRedeem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var redeemRequest NewRedeemRequest

	if err := json.NewDecoder(r.Body).Decode(&redeemRequest); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		http.Error(w, "bad cookie", http.StatusBadRequest)
	}

	err = auth.VerifyOTP(requestorRollno, redeemRequest.Otp)
	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Success:      false,
			ErrorMessage: "could not successfully verify otp",
		})
		return
	}

	err = account.NewRedeem(requestorRollno, redeemRequest.NumCoins, redeemRequest.Item)

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

func AcceptRedeem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var redeemRequest UpdateRedeemRequest

	if err := json.NewDecoder(r.Body).Decode(&redeemRequest); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		http.Error(w, "bad cookie", http.StatusBadRequest)
	}

	requestorRole := account.GetAccountRoleByRollno(requestorRollno)

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead) {
		http.Error(w, "you are not authorized to accept redeem requests", http.StatusUnauthorized)
		return
	}

	err = account.AcceptRedeem(redeemRequest.RedeemId, requestorRollno)

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

func RejectRedeem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var redeemRequest UpdateRedeemRequest

	if err := json.NewDecoder(r.Body).Decode(&redeemRequest); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		http.Error(w, "bad cookie", http.StatusBadRequest)
	}

	requestorRole := account.GetAccountRoleByRollno(requestorRollno)

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead) {
		http.Error(w, "you are not authorized to reject redeem requests", http.StatusUnauthorized)
		return
	}

	err = account.RejectRedeem(redeemRequest.RedeemId, requestorRollno)

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

func RedeemListByRollno(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	queriedRollno := r.URL.Query().Get("rollno")

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		http.Error(w, "bad cookie", http.StatusBadRequest)
	}

	requestorRole := account.GetAccountRoleByRollno(requestorRollno)

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead || requestorRollno != queriedRollno) {
		http.Error(w, "you are not authorized to view the requested redeem requests", http.StatusUnauthorized)
		return
	}

	redeemList, err := account.GetRedeemListByRollno(queriedRollno)

	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(&RedeemListResponse{
		Response: Response{
			Success: true,
		},
		RedeemList: redeemList,
	})
}
