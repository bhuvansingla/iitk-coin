package server

import (
	"net/http"

	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/errors"
	"github.com/bhuvansingla/iitk-coin/handlers"
)

func setRoutes() {
	http.HandleFunc("/", CORS(handlers.Index))
	http.HandleFunc("/auth/login", CORS(errors.Handler(handlers.Login)))
	http.HandleFunc("/auth/signup", CORS(errors.Handler(handlers.Signup)))
	http.HandleFunc("/auth/check", CORS(auth.IsAuthorized((errors.Handler(handlers.CheckLogin)))))
	http.HandleFunc("/auth/otp", CORS(errors.Handler(handlers.GenerateOtp)))
	http.HandleFunc("/auth/logout", CORS(errors.Handler(handlers.Logout)))

	http.HandleFunc("/user/name", CORS(auth.IsAuthorized(errors.Handler(handlers.GetNameByRollNo))))
	http.HandleFunc("/wallet/history", CORS(auth.IsAuthorized(errors.Handler(handlers.WalletHistory))))
	http.HandleFunc("/wallet/transfer/tax", CORS(auth.IsAuthorized(errors.Handler(handlers.TransferTax))))
	http.HandleFunc("/wallet/transfer", CORS(auth.IsAuthorized(errors.Handler(handlers.TransferCoins))))
	http.HandleFunc("/wallet/balance", CORS(auth.IsAuthorized(errors.Handler(handlers.GetCoinBalance))))
	http.HandleFunc("/wallet/redeem/new", CORS(auth.IsAuthorized(errors.Handler(handlers.NewRedeem))))

	http.HandleFunc("/wallet/add", CORS(auth.IsAuthorized(errors.Handler(handlers.RewardCoins))))
	http.HandleFunc("/wallet/redeem/accept", CORS(auth.IsAuthorized(errors.Handler(handlers.AcceptRedeem))))
	http.HandleFunc("/wallet/redeem/reject", CORS(auth.IsAuthorized(errors.Handler(handlers.RejectRedeem))))
	http.HandleFunc("/wallet/redeem/requests", CORS(auth.IsAuthorized(errors.Handler(handlers.RedeemListByRollNo))))
}
