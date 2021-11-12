package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type TransferCoinRequest struct {
	NumCoins       int64    `json:"numCoins"`
	ReceiverRollNo string `json:"receiverRollNo"`
	Remarks        string `json:"remarks"`
	Otp            string `json:"otp"`
}

type TransferTaxRequest struct {
	NumCoins       int64    `json:"numCoins"`
	ReceiverRollNo string `json:"receiverRollNo"`
}

type TransferTaxResponse struct {
	RollNo string `json:"rollNo"`
	Tax    int64    `json:"tax"`
}

type TransferCoinResponse struct {
	TxnID string `json:"id"`
}

func TransferCoins(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	var transferCoinRequest TransferCoinRequest

	if err := json.NewDecoder(r.Body).Decode(&transferCoinRequest); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error decoding request body")
	}

	requestorRollNo, err := auth.GetRollNoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "invalid cookie")
	}

	if err = auth.VerifyOTP(requestorRollNo, transferCoinRequest.Otp); err != nil {
		return err
	}

	if (requestorRollNo == transferCoinRequest.ReceiverRollNo) {
		return errors.NewHTTPError(nil, http.StatusBadRequest, "cannot transfer to self")
	}

	id, err := account.TransferCoins(requestorRollNo, transferCoinRequest.ReceiverRollNo, transferCoinRequest.NumCoins, transferCoinRequest.Remarks)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(&TransferCoinResponse{TxnID: id})
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}

func TransferTax(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	var transferTaxRequest TransferTaxRequest

	if err := json.NewDecoder(r.Body).Decode(&transferTaxRequest); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error decoding request body")
	}

	requestorRollNo, err := auth.GetRollNoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "invalid cookie")
	}

	if (requestorRollNo == transferTaxRequest.ReceiverRollNo) {
		return errors.NewHTTPError(nil, http.StatusBadRequest, "cannot transfer to self")
	}

	tax, err := account.CalculateTransferTax(requestorRollNo, transferTaxRequest.ReceiverRollNo, transferTaxRequest.NumCoins)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(&TransferTaxResponse{
			RollNo: requestorRollNo,
			Tax:    tax,
	})
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}
