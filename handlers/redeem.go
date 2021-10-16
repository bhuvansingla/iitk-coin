package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
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
	RedeemList []account.RedeemRequest `json:"redeemList"`
}

func NewRedeem(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	var redeemRequest NewRedeemRequest
	if err := json.NewDecoder(r.Body).Decode(&redeemRequest); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error decoding request body")
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "Invalid cookie")
	}

	if err = auth.VerifyOTP(requestorRollno, redeemRequest.Otp); err != nil {
		return err
	}

	if err = account.NewRedeem(requestorRollno, redeemRequest.NumCoins, redeemRequest.Item); err != nil {
		return err
	}

	return nil
}

func AcceptRedeem(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	var redeemRequest UpdateRedeemRequest
	if err := json.NewDecoder(r.Body).Decode(&redeemRequest); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error decoding request body")

	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "Invalid cookie")
	}

	requestorRole, err := account.GetAccountRoleByRollno(requestorRollno)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error when getting account role")
	}

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead) {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "you are not authorized to accept redeem requests")
	}

	if err = account.AcceptRedeem(redeemRequest.RedeemId, requestorRollno); err != nil {
		return err
	}

	return nil
}

func RejectRedeem(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	var redeemRequest UpdateRedeemRequest

	if err := json.NewDecoder(r.Body).Decode(&redeemRequest); err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error decoding request body")
	}

	requestorRollno, err := auth.GetRollnoFromRequest(r)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "Invalid cookie")
	}

	requestorRole, err := account.GetAccountRoleByRollno(requestorRollno)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error when getting account role")
	}

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead) {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "you are not authorized to reject redeem requests")
	}

	if err = account.RejectRedeem(redeemRequest.RedeemId, requestorRollno); err != nil {
		return err
	}

	return nil
}

func RedeemListByRollno(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	queriedRollno := r.URL.Query().Get("rollno")

	requestorRollno, err := auth.GetRollnoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "Invalid cookie")
	}

	requestorRole, err := account.GetAccountRoleByRollno(requestorRollno)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "error when getting account role")
	}

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead || requestorRollno != queriedRollno) {
		return errors.NewHTTPError(err, http.StatusUnauthorized, "you are not authorized to view the requested redeem requests")
	}

	redeemList, err := account.GetRedeemListByRollno(queriedRollno)

	if err != nil {
		return err
	}

	json.NewEncoder(w).Encode(&RedeemListResponse{
		RedeemList: redeemList,
	})

	return nil
}
