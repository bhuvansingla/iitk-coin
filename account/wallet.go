package account

import (
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
