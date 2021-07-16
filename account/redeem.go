package account

import (
	"errors"
	"time"

	"github.com/bhuvansingla/iitk-coin/database"
	log "github.com/sirupsen/logrus"
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
	stmt, err := database.DB.Prepare("INSERT INTO REDEEM_REQUEST (rollno,coins,time,status,item) VALUES (?,?,?,?,?)")
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
	err := database.DB.QueryRow("SELECT rollno, coins FROM REDEEM_REQUEST WHERE id=?", id).Scan(&redeemRequest)
	if err != nil {
		log.Error(err)
		return errors.New("internal server error")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return errors.New("transaction failed")
	}

	res, err := tx.Exec("UPDATE ACCOUNT SET coins = coins - ? WHERE rollno=? AND coins - ? >= 0", redeemRequest.NumCoins, redeemRequest.RollNo, redeemRequest.NumCoins)
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return errors.New("transaction failed")
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return errors.New("transaction failed")
	}

	if rowCnt == 0 {
		tx.Rollback()
		log.Error(err)
		return errors.New("transaction failed")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		log.Error(err)
		return errors.New("transaction failed")
	}

	return nil
}

func RejectRedeem(id int, adminRollno string) error {
	stmt, err := database.DB.Prepare("UPDATE REDEEM_REQUEST SET (status,actionByRollno) VALUES (?,?) WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(Rejected, adminRollno)
	if err != nil {
		return err
	}
	return nil
}

func GetRedeemListByRollno(rollno string) ([]RedeemRequest, error) {
	rows, err := database.DB.Query("SELECT * FROM REDEEM_REQUEST WHERE rollno=?", rollno)
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
