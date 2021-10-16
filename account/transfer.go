package account

import (
	"net/http"
	"time"

	"github.com/bhuvansingla/iitk-coin/database"
	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/spf13/viper"
)

func TransferCoins(fromRollno string, toRollno string, numCoins int, remarks string) error {

	err := validateCoinValue(numCoins)
	if err != nil {
		return err
	}

	userExistsFrom, err := UserExists(fromRollno)

	if err != nil {
		return err
	}

	userExistsTo, err := UserExists(toRollno)

	if err != nil {
		return err
	}

	if !userExistsFrom || !userExistsTo {
		return errors.NewHTTPError(nil, http.StatusBadRequest, "user account does not exist")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	res, err := tx.Exec("UPDATE ACCOUNT SET coins = coins - ? WHERE rollno = ? AND coins - ? >= 0 AND coins", numCoins, fromRollno, numCoins)
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

	limit := viper.GetInt("WALLET.UPPER_COIN_LIMIT")
	numCoinsToAdd := numCoins - calculateTax(fromRollno, toRollno, numCoins)

	res, err = tx.Exec("UPDATE ACCOUNT SET coins = coins + ? WHERE rollno=? AND coins + ? <= ?", numCoinsToAdd, toRollno, numCoinsToAdd, limit)
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
		return errors.NewHTTPError(nil, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	_, err = tx.Exec("INSERT INTO TRANSFER_HISTORY (fromRollno, toRollno, time, coins, tax, remarks) VALUES (?,?,?,?,?,?)", fromRollno, toRollno, time.Now(), numCoins, 0, remarks)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func calculateTax(rollno1 string, rollno2 string, numCoins int) (tax int) {
	if rollno1[:2] == rollno2[:2] {
		return (numCoins * viper.GetInt("TAX.INTER_BATCH") / 100)
	} else {
		return (numCoins * viper.GetInt("TAX.INTRA_BATCH") / 100)
	}
}
