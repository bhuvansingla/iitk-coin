package account

import (
	"database/sql"

	"github.com/bhuvansingla/iitk-coin/database"
)

type TransactionType string

const (
	REDEEM    TransactionType = "REDEEM"
	REWARD    TransactionType = "REWARD"
	TRANSFER  TransactionType = "TRANSFER"
)

type RedeemHistory struct {
	Type       TransactionType  `json:"type"`
	Time       int              `json:"timeStamp"`
	Id         string           `json:"txnID"`
	Amount     int              `json:"amount"`
	Item       string           `json:"item"`
	Status     RedeemStatus     `json:"status"`
	ActionByRollNo  string      `json:"actionByRollNo"`
	Name       string           `json:"name"`
}

type RewardHistory struct {
	Type    TransactionType  `json:"type"`
	Time    int	             `json:"timeStamp"`
	Id      string           `json:"txnID"`
	Amount  int	             `json:"amount"`
	Remarks string           `json:"remarks"`
}

type TransferHistory struct {
	Type        TransactionType `json:"type"`
	Time        int             `json:"timeStamp"`
	Id          string          `json:"txnID"`
	Amount      int             `json:"amount"`
	Tax         int             `json:"tax"`
	FromRollNo  string          `json:"fromRollNo"`
	ToRollNo    string          `json:"toRollNo"`
	Remarks     string          `json:"remarks"`
	Name        string          `json:"name"`
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
	SELECT history.*, a.name
	FROM (
		SELECT id,
			time,
			$2 AS type,
			fromRollNo,
			toRollNo,
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
			coins,
			NULL AS tax,
			NULL AS item,
			NULL AS status,
			NULL AS actionByRollNo,
			remarks
		FROM REWARD_HISTORY
		WHERE rollNo = $1
	) history
	LEFT JOIN ACCOUNT a
	ON (
		history.type = $2 AND a.rollNo = CASE WHEN history.fromRollNo = $1 THEN history.toRollNo ELSE history.fromRollNo END
		OR
		history.type = $3 AND a.rollNo = history.actionByRollNo
	)
	ORDER BY history.time DESC;`

	rows, err := database.DB.Query(queryString, rollNo, TRANSFER, REDEEM, REWARD)

	if err != nil {
		return nil, err
	}

	var history []interface{}

	for rows.Next() {
		var (
			id          int
			time        int
			txType      TransactionType
			fromRollNo  sql.NullString
			toRollNo    sql.NullString
			coins       sql.NullInt64
			tax         sql.NullInt64
			item        sql.NullString
			status      sql.NullString
			actionByRollNo sql.NullString
			remarks     sql.NullString
			name        sql.NullString
		)

		if err := rows.Scan(&id, &time, &txType, &fromRollNo, &toRollNo, &coins, &tax, &item, &status, &actionByRollNo, &remarks, &name); err != nil {
			return nil, err
		}

		var historyItem interface{}

		switch txType {
		case REDEEM:
			historyItem = RedeemHistory{
				Type: txType,
				Time: int(time),
				Id: formatTxnID(id, REDEEM),
				Amount: int(coins.Int64),
				Item: item.String,
				Status: RedeemStatus(status.String),
				ActionByRollNo: actionByRollNo.String,
				Name: name.String,
			}
		case REWARD:
			historyItem = RewardHistory{
				Type: txType,
				Time: time,
				Id: formatTxnID(id, REWARD),
				Amount: int(coins.Int64),
				Remarks: remarks.String,
			}
		case TRANSFER:
			historyItem = TransferHistory{
				Type: txType,
				Time: time,
				Id: formatTxnID(id, TRANSFER),
				Amount: int(coins.Int64),
				Tax: int(tax.Int64),
				FromRollNo: fromRollNo.String,
				ToRollNo: toRollNo.String,
				Remarks: remarks.String,
				Name: name.String,
			}
		}

		history = append(history, historyItem)
	}

	return history, nil
}
