package account

import (
	"database/sql"

	"github.com/bhuvansingla/iitk-coin/database"
)

func UpdateToken(token string, rollNo string) error {
	_, err := database.DB.Exec(("UPDATE REFRESH_TOKEN SET token = $1 WHERE rollNo = $2"), token, rollNo)
	return err
}

func DeleteToken(rollNo string) error {
	return UpdateToken("", rollNo)
}

func InvalidateAllTokens() error {
	_, err := database.DB.Exec(("UPDATE REFRESH_TOKEN SET token = $1"), "")
	return err
}

func GetToken(rollNo string) (string, error) {
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
