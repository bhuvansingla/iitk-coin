package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

var DB *sql.DB

func ConnectDB() error {
	name := viper.GetString("DATABASE.NAME")
	user := viper.GetString("DATABASE.USER")
	password := viper.GetString("DATABASE.PASS")
	dataSource := fmt.Sprintf("%s?_auth&_auth_user=%s&_auth_pass=%s", name, user, password)
	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func CreateTables() (err error) {
	err = createAccountTable()
	if err != nil {
		return
	}
	err = createOtpTable()
	if err != nil {
		return
	}
	err = createTransferHistoryTable()
	if err != nil {
		return
	}
	err = createRedeemRequestTable()
	if err != nil {
		return
	}
	return
}

func createAccountTable() (err error) {
	_, err = DB.Exec("create table if not exists ACCOUNT (rollno text, name text, password text, coins int, role int)")
	return
}

func createOtpTable() (err error) {
	_, err = DB.Exec("create table if not exists OTP (rollno text, otp text, created timestamp, used boolean)")
	return
}

func createTransferHistoryTable() (err error) {
	_, err = DB.Exec("create table if not exists TRANSFER_HISTORY (fromRollno text, toRollno text, time timestamp, coins int, tax int, remarks text)")
	return
}

func createRedeemRequestTable() (err error) {
	_, err = DB.Exec("create table if not exists REDEEM_REQUEST (id integer PRIMARY KEY AUTOINCREMENT NOT NULL, rollno text, coins int, time timestamp, item text, status text, actionByRollno text)")
	return
}
