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
	log "github.com/sirupsen/logrus"
)

func GenerateOtp(rollno string) error {

	if err := account.ValidateRollNo(rollno); err != nil {
		return err
	}

	validOtpExists, err := validOtpExists(rollno)
	if err != nil {
		return err
	}
	if validOtpExists {
		return errors.NewHTTPError(nil, http.StatusTooManyRequests, "OTP already sent. Please wait for some time.")
	}

	otp := util.RandomOTP()

	stmt, err := database.DB.Prepare("INSERT INTO OTP (rollno, otp, created, used) VALUES ($1,$2,$3,$4)")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(rollno, otp, time.Now().Unix(), 0); err != nil {
		return err
	}

	if err = mail.SendOTP(rollno, otp); err != nil {
		return err
	}

	log.Info(otp)

	return nil
}

func validOtpExists(rollno string) (bool, error) {
	createdBefore := time.Now().Add(-20 * time.Minute).Unix()

	row := database.DB.QueryRow("SELECT rollno FROM OTP WHERE rollno=$1 AND created > $2 AND used IS FALSE", rollno, createdBefore)
	var tempScan string
	err := row.Scan(&tempScan)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return true, err
	}
	return true, nil
}

func markOtpAsUsed(rollno string) error {
	_, err := database.DB.Exec("UPDATE OTP SET used=$1 WHERE rollno=$2", 1, rollno)
	if err != nil {
		return err
	}
	return nil
}

func VerifyOTP(rollno string, otp string) (err error) {
	createdBefore := time.Now().Add(-20 * time.Minute).Unix()

	row := database.DB.QueryRow("SELECT rollno FROM OTP WHERE rollno=$1 AND created > $2 AND otp=$3 AND used IS FALSE", rollno, createdBefore, otp)
	var tempScan string
	err = row.Scan(&tempScan)

	if err == sql.ErrNoRows {
		return errors.NewHTTPError(nil, http.StatusUnauthorized, "Invalid OTP")
	}
	if err != nil {
		return
	}
	
	err = markOtpAsUsed(rollno)
	return
}
