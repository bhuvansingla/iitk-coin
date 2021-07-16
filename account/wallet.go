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
	row := db.DB.QueryRow("SELECT coins FROM ACCOUNT WHERE rollno=?", rollno)
	var coins int
	if err := row.Scan(&coins); err != nil {
		log.Error("row scan failed")
		return 0, errors.New("internal server error")
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
