package wallet

import (
	"errors"
	"time"

	"github.com/bhuvansingla/iitk-coin/db"
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
	RollNo   string `field:"rollno"`
	NumCoins int    `field:"coins"`
}

func NewRedeem(rollno string, numCoins int, item string) error {
	stmt, err := db.DB.Prepare("INSERT INTO REDEEM_REQUEST (rollno,coins,time,status,item) VALUES (?,?,?,?,?)")
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
	err := db.DB.QueryRow("SELECT rollno, coins FROM REDEEM_REQUEST WHERE id=?", id).Scan(&redeemRequest)
	if err != nil {
		log.Error(err)
		return errors.New("internal server error")
	}

	tx, err := db.DB.Begin()
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
	stmt, err := db.DB.Prepare("UPDATE REDEEM_REQUEST SET (status,actionByRollno) VALUES (?,?) WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(Rejected, adminRollno)
	if err != nil {
		return err
	}
	return nil
}
