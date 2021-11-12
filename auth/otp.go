package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/database"
	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/bhuvansingla/iitk-coin/mail"
	"github.com/bhuvansingla/iitk-coin/util"

	"github.com/spf13/viper"
)

func GenerateOtp(rollNo string) error {

	if err := account.ValidateRollNo(rollNo); err != nil {
		return err
	}

	validOtpExists, err := validOtpExists(rollNo)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if validOtpExists {
		return errors.NewHTTPError(nil, http.StatusTooManyRequests, "please wait for some time. OTP already sent.")
	}

	otp := util.RandomOTP()

	stmt, err := database.DB.Prepare("INSERT INTO OTP (rollNo, otp, created, used) VALUES ($1,$2,$3,$4)")
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if _, err = stmt.Exec(rollNo, otp, time.Now().Unix(), 0); err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if err = mail.SendOTP(rollNo, otp); err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}

func validOtpExists(rollNo string) (bool, error) {
	newRequestWaitTime := viper.GetInt64("OTP.NEW_REQUEST_WAIT_TIME_IN_MIN")
	createdAfter := time.Now().Add(-time.Duration(newRequestWaitTime) * time.Minute).Unix()

	row := database.DB.QueryRow("SELECT rollNo FROM OTP WHERE rollNo=$1 AND created > $2 AND used IS FALSE", rollNo, createdAfter)
	var tempScan string

	err := row.Scan(&tempScan)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func markOtpAsUsed(rollNo string) error {
	_, err := database.DB.Exec("UPDATE OTP SET used=$1 WHERE rollNo=$2", 1, rollNo)
	if err != nil {
		return err
	}
	return nil
}

func VerifyOTP(rollNo string, otp string) (err error) {
	expiryPeriod := viper.GetInt64("OTP.EXPIRY_PERIOD_IN_MIN")
	createdAfter := time.Now().Add(-time.Duration(expiryPeriod) * time.Minute).Unix()

	row := database.DB.QueryRow("SELECT rollNo FROM OTP WHERE rollNo=$1 AND created > $2 AND otp=$3 AND used IS FALSE", rollNo, createdAfter, otp)
	var tempScan string
	err = row.Scan(&tempScan)

	if err == sql.ErrNoRows {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "invalid OTP")
	}
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	err = markOtpAsUsed(rollNo)
	if err != nil {
		return errors.NewHTTPError(err, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return
}
