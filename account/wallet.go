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

func GetCoinBalanceByRollNo(rollNo string) (int, error) {
	row := database.DB.QueryRow("SELECT coins FROM ACCOUNT WHERE rollNo=$1", rollNo)
	var coins int
	if err := row.Scan(&coins); err != nil {
		return 0, err
	}
	return coins, nil
}

func GetWalletHistoryByRollNo(rollNo string) ([]interface{}, error) {
	queryString := `
	SELECT history.*
	FROM (
		SELECT id,
			time,
			$2 AS type,
			fromRollNo,
			toRollNo,
			NULL AS rollNo,
			coins,
			tax,
			NULL AS item,
			NULL AS status,
			NULL AS actionByRollNo,
			remarks
		FROM TRANSFER_HISTORY
		WHERE toRollNo = $1 OR fromRollNo = $1
		UNION
		SELECT id,
			time,
			$3 AS type,
			NULL AS fromRollNo,
			NULL AS toRollNo,
			rollNo,
			coins,
			NULL AS tax,
			item,
			status,
			actionByRollNo,
			NULL AS remarks
		FROM REDEEM_REQUEST
		WHERE rollNo = $1
		UNION
		SELECT id,
			time,
			$4 AS type,
			NULL AS fromRollNo,
			NULL AS toRollNo,
			rollNo,
			coins,
			NULL AS tax,
			NULL AS item,
			NULL AS status,
			NULL AS actionByRollNo,
			remarks
		FROM REWARD_HISTORY
		WHERE rollNo = $1
	) history
	ORDER BY history.time DESC;`

	rows, err := database.DB.Query(queryString, rollNo, TRANSFER, REDEEM, REWARD)

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
			fromRollNo 	sql.NullString
			toRollNo	sql.NullString
			rollNo		sql.NullString
			coins		sql.NullInt64
			tax			sql.NullInt64
			item		sql.NullString
			status		sql.NullString
			actionByRollNo sql.NullString
			remarks		sql.NullString
		)
		
		if err := rows.Scan(&id, &time, &txType, &fromRollNo, &toRollNo, &rollNo, &coins, &tax, &item, &status, &actionByRollNo, &remarks); err != nil {
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
				FromRollNo: fromRollNo.String,
				ToRollNo: toRollNo.String,
				Remarks: remarks.String,
			}
		}

		history = append(history, historyItem)
	}

	return history, nil
}
