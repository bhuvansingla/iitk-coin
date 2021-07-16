package server

import (
	"net/http"

	"github.com/bhuvansingla/iitk-coin/auth"
	"github.com/bhuvansingla/iitk-coin/handlers"
)

func SetRoutes() {
	http.HandleFunc("/", CORS(handlers.Index))
	http.HandleFunc("/auth/login", CORS(handlers.Login))
	http.HandleFunc("/auth/signup", CORS(handlers.Signup))
	http.HandleFunc("/auth/check", CORS(auth.IsAuthorized(handlers.CheckLogin)))
	http.HandleFunc("/auth/otp", CORS(handlers.GenerateOtp))
	http.HandleFunc("/auth/logout", CORS(handlers.Logout))
	http.HandleFunc("/wallet/add", CORS(auth.IsAuthorized(handlers.RewardCoins)))
	http.HandleFunc("/wallet/transfer", CORS(auth.IsAuthorized(handlers.TransferCoins)))
	http.HandleFunc("/wallet/balance", CORS(auth.IsAuthorized(handlers.GetCoinBalance)))
	http.HandleFunc("/wallet/redeem/new", CORS(auth.IsAuthorized(handlers.NewRedeem)))
	http.HandleFunc("/wallet/redeem/accept", CORS(auth.IsAuthorized(handlers.AcceptRedeem)))
	http.HandleFunc("/wallet/redeem/reject", CORS(auth.IsAuthorized(handlers.RejectRedeem)))
}
