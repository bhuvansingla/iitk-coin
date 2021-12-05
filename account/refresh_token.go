package account

import (
	"database/sql"

	"github.com/bhuvansingla/iitk-coin/database"
)

func UpdateRefreshToken(token string, rollNo string) error {
	_, err := database.DB.Exec(("UPDATE REFRESH_TOKEN SET token = $1 WHERE rollNo = $2"), token, rollNo)
	return err
}

func DeleteRefreshToken(rollNo string) error {
	return UpdateRefreshToken("", rollNo)
}

func InvalidateAllRefreshTokens() error {
	_, err := database.DB.Exec(("UPDATE REFRESH_TOKEN SET token = $1"), "")
	return err
}

func GetRefreshToken(rollNo string) (string, error) {
	var token string
	err := database.DB.QueryRow(("SELECT token FROM REFRESH_TOKEN WHERE rollNo = $1"), rollNo).Scan(&token)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return token, nil
}
