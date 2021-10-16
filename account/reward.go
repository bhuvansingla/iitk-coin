package account

import (
	"net/http"

	"github.com/bhuvansingla/iitk-coin/database"
	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/spf13/viper"
)

func AddCoins(rollno string, coins int) error {

	if err := validateCoinValue(coins); err != nil {
		return err
	}

	userExists, err := UserExists(rollno)

	if err != nil {
		return err
	}

	if !userExists {
		return errors.NewHTTPError(nil, http.StatusBadRequest, "user account does not exist")
	}

	tx, err := database.DB.Begin()

	if err != nil {
		tx.Rollback()
		return err

	}

	limit := viper.GetInt("WALLET.UPPER_COIN_LIMIT")

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
		return errors.NewHTTPError(nil, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
