package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/bhuvansingla/iitk-coin/account"
	"github.com/bhuvansingla/iitk-coin/database"
	"github.com/bhuvansingla/iitk-coin/mail"
	"github.com/bhuvansingla/iitk-coin/util"
	log "github.com/sirupsen/logrus"
)

func GenerateOtp(rollno string) (string, error) {

	if err := account.ValidateRollNo(rollno); err != nil {
		return "", err
	}

	validOtpExists, err := validOtpExists(rollno)
	if err != nil {
		log.Error(err)
		return "", errors.New("internal server error")
	}
	if validOtpExists {
		return "", errors.New("otp exists already")
	}

	otp := util.RandomOTP()

	stmt, err := database.DB.Prepare("INSERT INTO OTP (rollno, otp, created, used) VALUES (?,?,?,?)")

	if err != nil {
		log.Error(err)
		return "", errors.New("internal server error")
	}
	_, err = stmt.Exec(rollno, otp, time.Now(), 0)
	if err != nil {
		log.Error(err)
		return "", errors.New("internal server error")
	}

	err = mail.SendOTP(rollno, otp)
	if err != nil {
		log.Error(err)
		return "", errors.New("internal server error")
	}

	log.Info(otp)
	return otp, nil
}

func validOtpExists(rollno string) (bool, error) {
	row := database.DB.QueryRow("SELECT rollno FROM OTP WHERE rollno=? AND created > datetime('now',  '-20 minute' , 'localtime') AND used IS FALSE", rollno)
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
	_, err := database.DB.Exec("UPDATE OTP SET used=? WHERE rollno=?", 1, rollno)
	if err != nil {
		return err
	}
	return nil
}

func VerifyOTP(rollno string, otp string) (err error) {
	row := database.DB.QueryRow("SELECT rollno FROM OTP WHERE rollno=? AND otp=? AND created > datetime('now',  '-20 minute' , 'localtime') AND used IS FALSE", rollno, otp)
	var tempScan string
	err = row.Scan(&tempScan)
	if err != nil {
		return
	}
	err = markOtpAsUsed(rollno)
	return
}
