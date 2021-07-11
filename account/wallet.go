package account

import (
	"database/sql"
	"errors"
	"time"

	"github.com/bhuvansingla/iitk-coin/db"
	log "github.com/sirupsen/logrus"
)

func GetCoinBalanceByRollno(rollno string) (int, error) {
	if !UserExists(rollno) {
		return 0, errors.New("user account does not exist")
	}
	row := db.DB.QueryRow("SELECT coins FROM ACCOUNT WHERE rollno=?", rollno)
	var coins int
	err := row.Scan(&coins)
	if err != nil {
		return 0, err
	}
	return coins, nil
}

func UpdateCoinBalanceByRollno(tx *sql.Tx, rollno string, coins int) error {
	_, err := tx.Exec("UPDATE ACCOUNT SET coins=? WHERE rollno=?", coins, rollno)
	if err != nil {
		return err
	}
	return nil
}

func AddCoins(rollno string, coins int) error {
	log.SetLevel(log.DebugLevel)

	err := validateCoinValue(coins)
	if err != nil {
		return err
	}

	if !UserExists(rollno) {
		return errors.New("user account does not exist")
	}

	tx, err := db.DB.Begin()

	if err != nil {
		tx.Rollback()
		return err
	}

	limit := 1000

	res, err := tx.Exec("UPDATE ACCOUNT SET coins = coins + ? WHERE rollno=? AND coins + ? <= ?", coins, rollno, coins, limit)
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
		return errors.New("upper limit exceeded")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

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

func validateCoinValue(coins int) error {
	if coins <= 0 {
		return errors.New("invalid coin value")
	}
	return nil
}
