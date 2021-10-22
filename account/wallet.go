package account

import (
	"database/sql"
	"time"

	"github.com/bhuvansingla/iitk-coin/database"
)

type transactionType string

const (
	REWARD 		transactionType = "REWARD"
	REDEEM 		transactionType = "REDEEM"
	TRANSFER 	transactionType = "TRANSFER"
)

type rewardHistory struct {
	Type	transactionType 	`json:"Type"`
	Time	int64				`json:"TimeStamp"`
	Id		string				`json:"TxnID"`
	Amount	int64				`json:"Amount"`
	Remarks string				`json:"Remarks"`
}

type transferHistory struct {
	Type	transactionType 	`json:"Type"`
	Time	int64				`json:"TimeStamp"`
	Id		string				`json:"TxnID"`
	Amount	int64				`json:"Amount"`
	Tax		int64				`json:"Tax"`
	FromRollNo string			`json:"FromRollNo"`
	ToRollNo string				`json:"ToRollNo"`
	Remarks string				`json:"Remarks"`
}

type redeemHistory struct {
	Type	transactionType 	`json:"Type"`
	Time	int64				`json:"TimeStamp"`
	Id		string				`json:"TxnID"`
	Amount	int64				`json:"Amount"`
	Remarks string				`json:"Remarks"`
	Status	RedeemStatus		`json:"Status"`
}

func GetCoinBalanceByRollno(rollno string) (int, error) {
	row := database.DB.QueryRow("SELECT coins FROM ACCOUNT WHERE rollno=?", rollno)
	var coins int
	if err := row.Scan(&coins); err != nil {
		return 0, err
	}
	return coins, nil
}

func GetWalletHistoryByRollno(rollno string) ([]interface{}, error) {
	queryString := `
	SELECT history.*
	FROM (
		SELECT id,
			time,
			"TRANSFER" AS type,
			fromRollno,
			toRollno,
			NULL AS rollno,
			coins,
			tax,
			NULL AS item,
			NULL AS status,
			NULL AS actionByRollno,
			remarks
		FROM TRANSFER_HISTORY
		WHERE toRollno = $1 OR fromRollno = $1
		UNION
		SELECT id,
			time,
			"REDEEM" AS type,
			NULL AS fromRollno,
			NULL AS toRollno,
			rollno,
			coins,
			NULL AS tax,
			item,
			status,
			actionByRollno,
			NULL AS remarks
		FROM REDEEM_REQUEST
		WHERE rollno = $1
		UNION
		SELECT id,
			time,
			"REWARD" AS type,
			NULL AS fromRollno,
			NULL AS toRollno,
			rollno,
			coins,
			NULL AS tax,
			NULL AS item,
			NULL AS status,
			NULL AS actionByRollno,
			remarks
		FROM REWARD_HISTORY
		WHERE rollno = $1
	) history
	ORDER BY history.time;`

	rows, err := database.DB.Query(queryString, rollno)

	if err != nil {
		return nil, err
	}

	var history []interface{}

	for rows.Next() {
		var (
			id 			string
			time 		time.Time
			txType 		transactionType
			fromRollno 	sql.NullString
			toRollno	sql.NullString
			rollno		sql.NullString
			coins		sql.NullInt64
			tax			sql.NullInt64
			item		sql.NullString
			status		sql.NullString
			actionByRollno sql.NullString
			remarks		sql.NullString
		)
		
		if err := rows.Scan(&id, &time, &txType, &fromRollno, &toRollno, &rollno, &coins, &tax, &item, &status, &actionByRollno, &remarks); err != nil {
			return nil, err
		}

		var historyItem interface{}
		switch txType {
		case REWARD:
			historyItem = rewardHistory{
				Type: txType,
				Time: time.Unix(),
				Id: "RWD" + id,
				Amount: coins.Int64,
				Remarks: remarks.String,
			}
		case REDEEM:
			historyItem = redeemHistory{
				Type: txType,
				Time: time.Unix(),
				Id: "REDM" + id,
				Amount: coins.Int64,
				Remarks: remarks.String,
				Status: RedeemStatus(status.String),
			}
		case TRANSFER:
			historyItem = transferHistory{
				Type: txType,
				Time: time.Unix(),
				Id: "TRNS" + id,
				Amount: coins.Int64,
				Tax: tax.Int64,
				FromRollNo: fromRollno.String,
				ToRollNo: toRollno.String,
				Remarks: remarks.String,
			}
		default:
			historyItem = "ERR"
		}

		history = append(history, historyItem)
	}

	return history, nil
}
