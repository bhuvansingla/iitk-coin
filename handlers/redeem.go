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
	ReceiverRollNo string `json:"receiverRollNo"`
	Item           string `json:"item"`
	Otp            string `json:"otp"`
}

type NewRedeemResponse struct {
	TxnId string `json:"id"`
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

	requestorRollNo, err := auth.GetRollNoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "invalid cookie")
	}

	if err = auth.VerifyOTP(requestorRollNo, redeemRequest.Otp); err != nil {
		return err
	}

	id, err := account.NewRedeem(requestorRollNo, redeemRequest.NumCoins, redeemRequest.Item)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(&NewRedeemResponse{TxnId: id})
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
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

	requestorRollNo, err := auth.GetRollNoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "invalid cookie")
	}

	requestorRole, err := account.GetAccountRoleByRollNo(requestorRollNo)
	if err != nil {
		return err
	}

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead) {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "you are not authorized to accept redeem requests")
	}

	if err = account.AcceptRedeem(redeemRequest.RedeemId, requestorRollNo); err != nil {
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

	requestorRollNo, err := auth.GetRollNoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "invalid cookie")
	}

	requestorRole, err := account.GetAccountRoleByRollNo(requestorRollNo)
	if err != nil {
		return err
	}

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead) {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "you are not authorized to reject redeem requests")
	}

	if err = account.RejectRedeem(redeemRequest.RedeemId, requestorRollNo); err != nil {
		return err
	}

	return nil
}

func RedeemListByRollNo(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	queriedRollNo := r.URL.Query().Get("rollNo")

	requestorRollNo, err := auth.GetRollNoFromRequest(r)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "invalid cookie")
	}

	requestorRole, err := account.GetAccountRoleByRollNo(requestorRollNo)
	if err != nil {
		return err
	}

	if !(requestorRole == account.GeneralSecretary || requestorRole == account.AssociateHead || requestorRollNo != queriedRollNo) {
		return errors.NewHTTPError(err, http.StatusUnauthorized, "you are not authorized to view the requested redeem requests")
	}

	redeemList, err := account.GetRedeemListByRollNo(queriedRollNo)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(&RedeemListResponse{
		RedeemList: redeemList,
	})
	if err != nil {
		return err
	}

	return nil
}
