package account

import (
	"errors"
	"time"

	"github.com/bhuvansingla/iitk-coin/db"
	log "github.com/sirupsen/logrus"
)

func TransferCoins(fromRollno string, toRollno string, numCoins int, remarks string) error {

	err := validateCoinValue(numCoins)
	if err != nil {
		return err
	}

	if !UserExists(fromRollno) || !UserExists(toRollno) {
		return errors.New("user account does not exist")
	}

	tx, err := db.DB.Begin()
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return errors.New("transaction failed")
	}

	res, err := tx.Exec("UPDATE ACCOUNT SET coins = coins - ? WHERE rollno = ? AND coins - ? >= 0 AND coins", numCoins, fromRollno, numCoins)
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
		log.Error(errors.New("no row changed"))
		return errors.New("transaction falied")
	}

	limit := 1000
	res, err = tx.Exec("UPDATE ACCOUNT SET coins = coins + ? WHERE rollno=? AND coins + ? <= ?", numCoins, toRollno, numCoins, limit)
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return errors.New("transaction failed")
	}

	rowCnt, err = res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowCnt == 0 {
		tx.Rollback()
		log.Error(errors.New("no row changed"))
		return errors.New("transaction falied")
	}
	log.Info(fromRollno, toRollno, time.Now(), numCoins, 0, remarks)
	_, err = tx.Exec("INSERT INTO TRANSFER_HISTORY (fromRollno, toRollno, time, coins, tax, remarks) VALUES (?,?,?,?,?,?)", fromRollno, toRollno, time.Now(), numCoins, 0, remarks)
	if err != nil {
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