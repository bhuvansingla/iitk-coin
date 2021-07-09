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
	http.Handle("/auth/check", auth.IsAuthorized(handlers.CheckUserIsLoggedIn))
	http.HandleFunc("/auth/otp", handlers.GenerateOtp)
	http.Handle("/coins/add", auth.IsAuthorized(handlers.AddCoins))
	http.Handle("/coins/transfer", auth.IsAuthorized(handlers.TransferCoins))
	http.Handle("/coins/balance", auth.IsAuthorized(handlers.GetCoinBalance))
}
