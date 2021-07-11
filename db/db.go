package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDB() error {
	db, err := sql.Open("sqlite3", "iitkcoin.db")
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
