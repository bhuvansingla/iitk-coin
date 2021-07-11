package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/wallet"
)

type TransferCoinRequest struct {
	NumCoins       int    `json:"numCoins"`
	ReceiverRollno string `json:"receiverRollno"`
	Remarks        string `json:"remarks"`
	Otp            string `json:"otp"`
}

func TransferCoins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var transferCoinRequest TransferCoinRequest

	if err := json.NewDecoder(r.Body).Decode(&transferCoinRequest); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		http.Error(w, "bad cookie", http.StatusBadRequest)
	}

	err = auth.VerifyOTP(requestorRollno, transferCoinRequest.Otp)
	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Success:      false,
			ErrorMessage: "could not successfully verify otp",
		})
		return
	}

	err = wallet.TransferCoins(requestorRollno, transferCoinRequest.ReceiverRollno, transferCoinRequest.NumCoins, transferCoinRequest.Remarks)

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
