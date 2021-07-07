package auth

import (
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/bhuvansingla/iitk-coin/db"
	"github.com/sirupsen/logrus"
)

func GenerateOtp(rollno string) (string, error) {

	validOtpExists, err := ValidOtpExists(rollno)
	if err != nil {
		return "", err
	}
	if validOtpExists {
		return "", errors.New("otp exists already")
	}

	stmt, err := db.DB.Prepare("INSERT INTO OTPs (rollno, otp, created, used) VALUES (?,?,?,?)")

	if err != nil {
		return "", err
	}

	otp := randomString()
	_, err = stmt.Exec(rollno, otp, time.Now(), 0)

	if err != nil {
		return "", err
	}
	logrus.Info(otp)
	return otp, nil
}

func ValidOtpExists(rollno string) (bool, error) {
	row := db.DB.QueryRow("SELECT rollno FROM OTPs WHERE rollno=? AND created > datetime('now',  '-20 minute' , 'localtime') AND used IS FALSE", rollno)
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

func MarkOtpAsUsed(rollno string) error {
	_, err := db.DB.Exec("UPDATE OTPs SET used=? WHERE rollno=?", 1, rollno)
	if err != nil {
		return err
	}
	return nil
}

func VerifyOTP(rollno string, otp string) (bool, error) {
	row := db.DB.QueryRow("SELECT rollno FROM OTPs WHERE rollno=? AND otp=? AND created > datetime('now',  '-20 minute' , 'localtime') AND used IS FALSE", rollno, otp)
	var tempScan string
	err := row.Scan(&tempScan)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return true, err
	}
	err = MarkOtpAsUsed(rollno)
	return true, err
}

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomString() string {
	return StringWithCharset(10, charset)
}
