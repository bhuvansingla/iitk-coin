package account

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/bhuvansingla/iitk-coin/database"
	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/spf13/viper"
)

type RedeemStatus string

const (
	Pending   RedeemStatus = "PENDING"
	Cancelled RedeemStatus = "CANCELLED"
	Approved  RedeemStatus = "APPROVED"
	Rejected  RedeemStatus = "REJECTED"
)

type RedeemRequest struct {
	RollNo   string       `field:"rollNo"`
	Id       string       `field:"id"`
	NumCoins int          `field:"coins"`
	Time     int          `field:"time"`
	Item     string       `field:"item"`
	Status   RedeemStatus `field:"status"`
	ActionByRollNo string `field:"actionByRollNo"`
}

func NewRedeem(rollNo string, numCoins int, item string) (string, error) {
	var (
		redeemSuffix = viper.GetString("TXNID.REDEEM_SUFFIX")
		txnIDPadding = viper.GetInt("TXNID.PADDING")
		id int
	)

	stmt, err := database.DB.Prepare("INSERT INTO REDEEM_REQUEST (rollNo,coins,time,status,item) VALUES ($1,$2,$3,$4,$5) RETURNING id")
	if err != nil {
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	err = stmt.QueryRow(rollNo, numCoins, time.Now().Unix(), Pending, item).Scan(&id)
	if err != nil {
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	
	return fmt.Sprintf("%s%0*d", redeemSuffix, txnIDPadding, id),  nil
}

func AcceptRedeem(id int, adminRollNo string) error {

	var redeemRequest RedeemRequest
	err := database.DB.QueryRow("SELECT rollNo, coins FROM REDEEM_REQUEST WHERE id=$1", id).Scan(&redeemRequest.RollNo, &redeemRequest.NumCoins)
	if err == sql.ErrNoRows {
		return errors.NewHTTPError(err, http.StatusNotFound, "invalid redeem ID")
	} else if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	res, err := tx.Exec("UPDATE ACCOUNT SET coins = coins - $1 WHERE rollNo=$2 AND coins >= $1", redeemRequest.NumCoins, redeemRequest.RollNo)
	if err != nil {
		tx.Rollback()
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if rowCnt == 0 {
		tx.Rollback()
		return errors.NewHTTPError(nil, http.StatusBadRequest, "insufficient wallet balance")
	}

	_, err = tx.Exec("UPDATE REDEEM_REQUEST SET status=$1, actionByRollNo=$2 WHERE id=$3", Approved, adminRollNo, id)
	if err != nil {
		tx.Rollback()
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}

func RejectRedeem(id int, adminRollNo string) error {
	stmt, err := database.DB.Prepare("UPDATE REDEEM_REQUEST SET status=$1, actionByRollNo=$2 WHERE id=$3")
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	res, err := stmt.Exec(Rejected, adminRollNo, id)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if rowCnt == 0 {
		return errors.NewHTTPError(nil, http.StatusBadRequest, "invalid redeem ID")
	}

	return nil
}

func GetRedeemListByRollNo(rollNo string) ([]RedeemRequest, error) {
	rows, err := database.DB.Query("SELECT * FROM REDEEM_REQUEST WHERE rollNo=$1", rollNo)
	if err != nil {
		return nil, errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	defer rows.Close()

	var redeemRequests []RedeemRequest
	for rows.Next() {
		var redeemRequest RedeemRequest
		var adminRollNo sql.NullString
		err := rows.Scan(&redeemRequest.Id, &redeemRequest.RollNo, &redeemRequest.NumCoins, &redeemRequest.Time, &redeemRequest.Item, &redeemRequest.Status, &adminRollNo)
		if err != nil {
			return nil, errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		redeemRequest.ActionByRollNo = adminRollNo.String
		redeemRequests = append(redeemRequests, redeemRequest)
	}
	return redeemRequests, nil
}
