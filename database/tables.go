package database

import (
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func createTables() (err error) {
	err = createAccountTable()
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = createOtpTable()
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = createTransferHistoryTable()
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = createRedeemRequestTable()
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = createRewardHistoryTable()
	if err != nil {
		log.Error(err.Error())
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
	_, err = DB.Exec("create table if not exists TRANSFER_HISTORY (id integer PRIMARY KEY AUTOINCREMENT NOT NULL, fromRollno text, toRollno text, time timestamp, coins int, tax int, remarks text)")
	return
}

func createRedeemRequestTable() (err error) {
	_, err = DB.Exec("create table if not exists REDEEM_REQUEST (id integer PRIMARY KEY AUTOINCREMENT NOT NULL, rollno text, coins int, time timestamp, item text, status text, actionByRollno text)")
	return
}

func createRewardHistoryTable() (err error) {
	_, err = DB.Exec("create table if not exists REWARD_HISTORY (id integer PRIMARY KEY AUTOINCREMENT NOT NULL, rollno text, coins int, time timestamp, remarks text)")
	return
}
