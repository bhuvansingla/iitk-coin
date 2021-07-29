package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type TransferCoinRequest struct {
	NumCoins       int    `json:"numCoins"`
	ReceiverRollno string `json:"receiverRollno"`
	Remarks        string `json:"remarks"`
	Otp            string `json:"otp"`
}

func TransferCoins(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	var transferCoinRequest TransferCoinRequest

	if err := json.NewDecoder(r.Body).Decode(&transferCoinRequest); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error decoding request body")
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "Invalid cookie")
	}

	if err = auth.VerifyOTP(requestorRollno, transferCoinRequest.Otp); err != nil {
		return err
	}

	if err = account.TransferCoins(requestorRollno, transferCoinRequest.ReceiverRollno, transferCoinRequest.NumCoins, transferCoinRequest.Remarks); err != nil {
		return err
	}

	return nil
}
