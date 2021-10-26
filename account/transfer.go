package account

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bhuvansingla/iitk-coin/database"
	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/spf13/viper"
)

func TransferCoins(fromRollno string, toRollno string, numCoins int, remarks string) (string, error) {

	err := validateCoinValue(numCoins)
	if err != nil {
		return "", err
	}

	userExistsFrom, err := UserExists(fromRollno)

	if err != nil {
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	userExistsTo, err := UserExists(toRollno)

	if err != nil {
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if !userExistsFrom || !userExistsTo {
		return "", errors.NewHTTPError(nil, http.StatusBadRequest, "user account does not exist")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return "", err
	}

	res, err := tx.Exec("UPDATE ACCOUNT SET coins = coins - $1 WHERE rollno = $2 AND coins - $1 >= 0 AND coins", numCoins, fromRollno)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if rowCnt == 0 {
		tx.Rollback()
		return "", errors.NewHTTPError(nil, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	limit := viper.GetInt("WALLET.UPPER_COIN_LIMIT")
	tax, err := CalculateTransferTax(fromRollno, toRollno, numCoins)
	if err != nil {
		tx.Rollback()
		return err
	}

	numCoinsToAdd := numCoins - tax

	res, err = tx.Exec("UPDATE ACCOUNT SET coins = coins + $1 WHERE rollno=$2 AND coins + $1 <= $3", numCoinsToAdd, toRollno, limit)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	rowCnt, err = res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if rowCnt == 0 {
		tx.Rollback()
		return "", errors.NewHTTPError(nil, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	stmt, err := tx.Prepare("INSERT INTO TRANSFER_HISTORY (fromRollno, toRollno, time, coins, tax, remarks) VALUES ($1, $2, $3, $4, $5, $6)  RETURNING id")


	if err != nil {
		tx.Rollback()
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	var (
		transferSuffix = viper.GetString("TXNID.TRANSFER_SUFFIX")
		txnIDPadding = viper.GetInt("TXNID.PADDING")
		id int
	)

	err = stmt.QueryRow(fromRollno, toRollno, time.Now().Unix(), numCoins, numCoins - numCoinsToAdd, remarks).Scan(&id);
	
	if err != nil {
		tx.Rollback()
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	if err != nil {
		tx.Rollback()
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return "", err
	}

	return fmt.Sprintf("%s%0*d", transferSuffix, txnIDPadding, id), nil
}

func CalculateTransferTax(fromRollno string, toRollno string, numCoins int) (int, error) {

	err := validateCoinValue(numCoins)
	if err != nil {
		return 0, err
	}

	userExistsFrom, err := UserExists(fromRollno)

	if err != nil {
		return 0, errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	userExistsTo, err := UserExists(toRollno)

	if err != nil {
		return 0, errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if !userExistsFrom || !userExistsTo {
		return 0, errors.NewHTTPError(nil, http.StatusBadRequest, "user account does not exist")
	}

	var tax int
	if fromRollno[:2] == toRollno[:2] {
		tax = (numCoins * viper.GetInt("TAX.INTER_BATCH") / 100)
	} else {
		tax = (numCoins * viper.GetInt("TAX.INTRA_BATCH") / 100)
	}

	return tax, nil
}
