package account

import (
	"errors"
)

func validateCoinValue(coins int) error {
	if coins <= 0 {
		return errors.New("invalid coin value")
	}
	return nil
}
