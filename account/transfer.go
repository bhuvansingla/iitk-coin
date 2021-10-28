package account

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bhuvansingla/iitk-coin/database"
	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/spf13/viper"
)

func TransferCoins(fromRollNo string, toRollNo string, numCoins int, remarks string) (string, error) {

	err := validateCoinValue(numCoins)
	if err != nil {
		return "", err
	}

	userExistsFrom, err := UserExists(fromRollNo)
	if err != nil {
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	userExistsTo, err := UserExists(toRollNo)
	if err != nil {
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if !userExistsFrom || !userExistsTo {
		return "", errors.NewHTTPError(nil, http.StatusBadRequest, "user account does not exist")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	res, err := tx.Exec("UPDATE ACCOUNT SET coins = coins - $1 WHERE rollNo = $2 AND coins - $1 >= 0", numCoins, fromRollNo)
	if err != nil {
		tx.Rollback()
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if rowCnt == 0 {
		tx.Rollback()
		return "", errors.NewHTTPError(nil, http.StatusBadRequest, "insufficient wallet balance")
	}

	limit := viper.GetInt("WALLET.UPPER_COIN_LIMIT")
	tax, err := CalculateTransferTax(fromRollNo, toRollNo, numCoins)
	if err != nil {
		tx.Rollback()
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	numCoinsToAdd := numCoins - tax

	res, err = tx.Exec("UPDATE ACCOUNT SET coins = coins + $1 WHERE rollNo=$2 AND coins + $1 <= $3", numCoinsToAdd, toRollNo, limit)
	if err != nil {
		tx.Rollback()
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	rowCnt, err = res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if rowCnt == 0 {
		tx.Rollback()
		return "", errors.NewHTTPError(nil, http.StatusBadRequest, "receiver wallet upper limit reached")
	}

	stmt, err := tx.Prepare("INSERT INTO TRANSFER_HISTORY (fromRollNo, toRollNo, time, coins, tax, remarks) VALUES ($1, $2, $3, $4, $5, $6)  RETURNING id")
	if err != nil {
		tx.Rollback()
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	var (
		transferSuffix = viper.GetString("TXNID.TRANSFER_SUFFIX")
		txnIDPadding = viper.GetInt("TXNID.PADDING")
		id int
	)

	err = stmt.QueryRow(fromRollNo, toRollNo, time.Now().Unix(), numCoins, tax, remarks).Scan(&id);
	
	if err != nil {
		tx.Rollback()
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return fmt.Sprintf("%s%0*d", transferSuffix, txnIDPadding, id), nil
}

func CalculateTransferTax(fromRollNo string, toRollNo string, numCoins int) (int, error) {

	err := validateCoinValue(numCoins)
	if err != nil {
		return 0, err
	}

	userExistsFrom, err := UserExists(fromRollNo)
	if err != nil {
		return 0, errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	userExistsTo, err := UserExists(toRollNo)
	if err != nil {
		return 0, errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if !userExistsFrom || !userExistsTo {
		return 0, errors.NewHTTPError(nil, http.StatusBadRequest, "user account does not exist")
	}

	var tax int
	if fromRollNo[:2] == toRollNo[:2] {
		tax = (numCoins * viper.GetInt("TAX.INTER_BATCH") / 100)
	} else {
		tax = (numCoins * viper.GetInt("TAX.INTRA_BATCH") / 100)
	}

	return tax, nil
}
