package server

import (
	"net/http"

	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/handlers"
)

func SetRoutes() {
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/auth/login", handlers.Login)
	http.HandleFunc("/auth/signup", handlers.Signup)
	http.Handle("/auth/check", auth.IsAuthorized(handlers.CheckLogin))
	http.HandleFunc("/auth/otp", handlers.GenerateOtp)
	http.Handle("/wallet/add", auth.IsAuthorized(handlers.RewardCoins))
	http.Handle("/wallet/transfer", auth.IsAuthorized(handlers.TransferCoins))
	http.Handle("/wallet/balance", auth.IsAuthorized(handlers.GetCoinBalance))
}
