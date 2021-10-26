package account

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bhuvansingla/iitk-coin/database"
	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/spf13/viper"
)

func AddCoins(rollno string, coins int, remarks string) (string, error) {

	if err := validateCoinValue(coins); err != nil {
		return "", err
	}

	userExists, err := UserExists(rollno)

	if err != nil {
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if !userExists {
		return "", errors.NewHTTPError(nil, http.StatusBadRequest, "user account does not exist")
	}

	tx, err := database.DB.Begin()

	if err != nil {
		tx.Rollback()
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	limit := viper.GetInt("WALLET.UPPER_COIN_LIMIT")

	res, err := tx.Exec("UPDATE ACCOUNT SET coins = coins + $1 WHERE rollno=$2 AND coins + $1 <= $3", coins, rollno, limit)

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
		return "", errors.NewHTTPError(nil, http.StatusBadRequest, "wallet upper limit reached")
	}

	stmt, err := tx.Prepare("INSERT INTO REWARD_HISTORY (rollno, coins, time, remarks) VALUES ($1, $2, $3, $4) RETURNING id")

	if err != nil {
		tx.Rollback()
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	var (
		rewardSuffix = viper.GetString("TXNID.REWARD_SUFFIX")
		txnIDPadding = viper.GetInt("TXNID.PADDING")
		id int
	)

	err = stmt.QueryRow(rollno, coins, time.Now().Unix(), remarks).Scan(&id);
	
	if err != nil {
		tx.Rollback()
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return "", errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return fmt.Sprintf("%s%0*d", rewardSuffix, txnIDPadding, id), nil
}
