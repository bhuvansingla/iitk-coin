package account

import (
	"database/sql"
	"errors"

	"github.com/bhuvansingla/iitk-coin/db"
	log "github.com/sirupsen/logrus"
)

func GetCoinBalanceByRollno(rollno string) (int, error) {
	if !UserExists(rollno) {
		return 0, errors.New("user account does not exist")
	}
	row := db.DB.QueryRow("SELECT coins FROM Account WHERE rollno=?", rollno)
	var coins int
	err := row.Scan(&coins)
	if err != nil {
		return 0, err
	}
	return coins, nil
}

func UpdateCoinBalanceByRollno(tx *sql.Tx, rollno string, coins int) error {
	_, err := tx.Exec("UPDATE Account SET coins=? WHERE rollno=?", coins, rollno)
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

	res, err := tx.Exec("UPDATE Account SET coins = coins + ? WHERE rollno=? AND coins + ? <= ?", coins, rollno, coins, limit)
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

	log.SetLevel(log.DebugLevel)

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
		return err
	}

	res, err := tx.Exec("UPDATE Account SET coins = coins - ? WHERE rollno = ? AND coins - ? >= 0 AND coins ", numCoins, fromRollno, numCoins)
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
		return errors.New("no row changed")
	}

	limit := 1000
	res, err = tx.Exec("UPDATE Account SET coins = coins + ? WHERE rollno=? AND coins + ? <= ?", numCoins, toRollno, numCoins, limit)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowCnt, err = res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowCnt == 0 {
		tx.Rollback()
		return errors.New("no row changed")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func validateCoinValue(coins int) error {
	if coins <= 0 {
		return errors.New("invalid coin value")
	}
	return nil
}
