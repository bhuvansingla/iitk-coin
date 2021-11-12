package account

import (
	"fmt"
	"net/http"

	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/spf13/viper"
)

func validateCoinValue(coins int64) error {
	if coins <= 0 {
		return errors.NewHTTPError(nil, http.StatusBadRequest, "coin value must be greater than 0")
	}
	return nil
}

func formatTxnID(id int64, txnType TransactionType) string {
	suffix := ""
	switch txnType {
		case REWARD:
			suffix = viper.GetString("TXNID.REWARD_SUFFIX")
		case REDEEM:
			suffix = viper.GetString("TXNID.REDEEM_SUFFIX")
		case TRANSFER:
			suffix = viper.GetString("TXNID.TRANSFER_SUFFIX")
		default:
			suffix = "TXN"
	}
	txnIDPadding := viper.GetInt64("TXNID.PADDING")

	return fmt.Sprintf("%s%0*d", suffix, txnIDPadding, id)
}
