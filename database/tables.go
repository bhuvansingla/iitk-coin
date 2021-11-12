package database

import (
	_ "github.com/lib/pq"
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
	err = createRefreshTokenTable()
	if err != nil {
		log.Error(err.Error())
		return
	}
	return
}

func createAccountTable() (err error) {
	_, err = DB.Exec("create table if not exists ACCOUNT (rollNo text PRIMARY KEY NOT NULL, name text, password text, coins int, role int)")
	return
}

func createOtpTable() (err error) {
	_, err = DB.Exec("create table if not exists OTP (rollNo text, otp text, created NUMERIC, used boolean)")
	return
}

func createTransferHistoryTable() (err error) {
	_, err = DB.Exec("create table if not exists TRANSFER_HISTORY (id SERIAL PRIMARY KEY NOT NULL, fromRollNo text, toRollNo text, time NUMERIC, coins int, tax int, remarks text)")
	return
}

func createRedeemRequestTable() (err error) {
	_, err = DB.Exec("create table if not exists REDEEM_REQUEST (id SERIAL PRIMARY KEY NOT NULL, rollNo text, coins int, time NUMERIC, item text, status text, actionByRollNo text)")
	return
}

func createRewardHistoryTable() (err error) {
	_, err = DB.Exec("create table if not exists REWARD_HISTORY (id SERIAL PRIMARY KEY NOT NULL, rollNo text, coins int, time NUMERIC, remarks text)")
	return
}

func createRefreshTokenTable() (err error) {
	_, err = DB.Exec("create table if not exists REFRESH_TOKEN (rollNo text PRIMARY KEY NOT NULL, token text)")
	return
}
