package account

import (
	"database/sql"

	"github.com/bhuvansingla/iitk-coin/database"
)

func GetCoinBalanceByRollno(rollno string) (int, error) {
	row := database.DB.QueryRow("SELECT coins FROM ACCOUNT WHERE rollno=?", rollno)
	var coins int
	if err := row.Scan(&coins); err != nil {
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
