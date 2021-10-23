package account

import (
	"net/http"
	"time"

	"github.com/bhuvansingla/iitk-coin/database"
	"github.com/bhuvansingla/iitk-coin/errors"
)

type RedeemStatus string

const (
	Pending   RedeemStatus = "PENDING"
	Cancelled RedeemStatus = "CANCELLED"
	Approved  RedeemStatus = "APPROVED"
	Rejected  RedeemStatus = "REJECTED"
)

type RedeemRequest struct {
	RollNo   string       `field:"rollno"`
	Id       string       `field:"id"`
	NumCoins int          `field:"coins"`
	Time     time.Time    `field:"time"`
	Item     string       `field:"item"`
	Status   RedeemStatus `field:"status"`
}

func NewRedeem(rollno string, numCoins int, item string) error {
	stmt, err := database.DB.Prepare("INSERT INTO REDEEM_REQUEST (rollno,coins,time,status,item) VALUES ($1,$2,$3,$4,$5)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(rollno, numCoins, time.Now(), Pending, item)
	if err != nil {
		return err
	}
	return nil
}

func AcceptRedeem(id int, adminRollno string) error {

	var redeemRequest RedeemRequest
	err := database.DB.QueryRow("SELECT rollno, coins FROM REDEEM_REQUEST WHERE id=$1", id).Scan(&redeemRequest)
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	res, err := tx.Exec("UPDATE ACCOUNT SET coins = coins - $1 WHERE rollno=$2 AND coins - $1 >= 0", redeemRequest.NumCoins, redeemRequest.RollNo)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowCnt == 0 {
		tx.Rollback()
		return errors.NewHTTPError(nil, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func RejectRedeem(id int, adminRollno string) error {
	stmt, err := database.DB.Prepare("UPDATE REDEEM_REQUEST SET (status,actionByRollno) VALUES ($1,$2) WHERE id=$3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(Rejected, adminRollno, id)
	if err != nil {
		return err
	}
	return nil
}

func GetRedeemListByRollno(rollno string) ([]RedeemRequest, error) {
	rows, err := database.DB.Query("SELECT * FROM REDEEM_REQUEST WHERE rollno=$1", rollno)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var redeemRequests []RedeemRequest
	for rows.Next() {
		var redeemRequest RedeemRequest
		err := rows.Scan(&redeemRequest.Id, &redeemRequest.RollNo, &redeemRequest.NumCoins, &redeemRequest.Time, &redeemRequest.Item, &redeemRequest.Status)
		if err != nil {
			return nil, err
		}
		redeemRequests = append(redeemRequests, redeemRequest)
	}
	return redeemRequests, nil
}
