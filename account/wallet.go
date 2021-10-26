package account

import (
	"database/sql"
	"fmt"

	"github.com/bhuvansingla/iitk-coin/database"
	"github.com/spf13/viper"
)

type TransactionType string

const (
	REWARD 		TransactionType = "REWARD"
	REDEEM 		TransactionType = "REDEEM"
	TRANSFER 	TransactionType = "TRANSFER"
)

type RewardHistory struct {
	Type	TransactionType 	`json:"type"`
	Time	int64				`json:"timeStamp"`
	Id		string				`json:"txnID"`
	Amount	int64				`json:"amount"`
	Remarks string				`json:"remarks"`
}

type TransferHistory struct {
	Type		TransactionType `json:"type"`
	Time		int64			`json:"timeStamp"`
	Id			string			`json:"txnID"`
	Amount		int64			`json:"amount"`
	Tax			int64			`json:"tax"`
	FromRollNo 	string			`json:"fromRollNo"`
	ToRollNo 	string			`json:"toRollNo"`
	Remarks 	string			`json:"remarks"`
}

type RedeemHistory struct {
	Type	TransactionType 	`json:"type"`
	Time	int64				`json:"timeStamp"`
	Id		string				`json:"txnID"`
	Amount	int64				`json:"amount"`
	Remarks string				`json:"remarks"`
	Status	RedeemStatus		`json:"status"`
}

func GetCoinBalanceByRollNo(rollno string) (int, error) {
	row := database.DB.QueryRow("SELECT coins FROM ACCOUNT WHERE rollno=$1", rollno)
	var coins int
	if err := row.Scan(&coins); err != nil {
		return 0, err
	}
	return coins, nil
}

func GetWalletHistoryByRollNo(rollno string) ([]interface{}, error) {
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
	ORDER BY history.time DESC;`

	rows, err := database.DB.Query(queryString, rollno)

	if err != nil {
		return nil, err
	}

	var history []interface{}

	var (
		redeemSuffix = viper.GetString("TXNID.REDEEM_SUFFIX")
		rewardSuffix = viper.GetString("TXNID.REWARD_SUFFIX")
		transferSuffix = viper.GetString("TXNID.TRANSFER_SUFFIX")
		txnIDPadding = viper.GetInt("TXNID.PADDING")
	)
	
	for rows.Next() {
		var (
			id 			int
			time 		int64
			txType 		TransactionType
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
			historyItem = RewardHistory{
				Type: txType,
				Time: time,
				Id: fmt.Sprintf("%s%0*d", rewardSuffix, txnIDPadding, id),
				Amount: coins.Int64,
				Remarks: remarks.String,
			}
		case REDEEM:
			historyItem = RedeemHistory{
				Type: txType,
				Time: time,
				Id: fmt.Sprintf("%s%0*d", redeemSuffix, txnIDPadding, id),
				Amount: coins.Int64,
				Remarks: remarks.String,
				Status: RedeemStatus(status.String),
			}
		case TRANSFER:
			historyItem = TransferHistory{
				Type: txType,
				Time: time,
				Id: fmt.Sprintf("%s%0*d", transferSuffix, txnIDPadding, id),
				Amount: coins.Int64,
				Tax: tax.Int64,
				FromRollNo: fromRollno.String,
				ToRollNo: toRollno.String,
				Remarks: remarks.String,
			}
		}

		history = append(history, historyItem)
	}

	return history, nil
}
