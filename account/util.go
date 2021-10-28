package account

import (
	"net/http"

	"github.com/bhuvansingla/iitk-coin/errors"
)

func validateCoinValue(coins int) error {
	if coins <= 0 {
		return errors.NewHTTPError(nil, http.StatusBadRequest, "coin value must be greater than 0")
	}
	return nil
}
