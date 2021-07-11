package wallet

import (
	"errors"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/db"
	log "github.com/sirupsen/logrus"
)

func AddCoins(rollno string, coins int) error {

	if err := validateCoinValue(coins); err != nil {
		return err
	}

	if !account.UserExists(rollno) {
		return errors.New("user account does not exist")
	}

	tx, err := db.DB.Begin()

	if err != nil {
		tx.Rollback()
		log.Error(err)
		return errors.New("transaction failed")

	}

	limit := 1000

	res, err := tx.Exec("UPDATE ACCOUNT SET coins = coins + ? WHERE rollno=? AND coins + ? <= ?", coins, rollno, coins, limit)
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
