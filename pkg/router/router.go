package router

import (
	"net/http"

	"github.com/bhuvansingla/iitk-coin/pkg/handlers"
	"github.com/bhuvansingla/iitk-coin/pkg/jwt"
)

func SetRoutes() {
	http.Handle("/", jwt.IsAuthorized(handlers.Index))
	http.HandleFunc("/login", handlers.Login)
	http.Handle("/checklogin", jwt.IsAuthorized(handlers.CheckUserIsLoggedIn))
	http.HandleFunc("/signup", handlers.Signup)
	http.Handle("/coins/add", jwt.IsAuthorized(handlers.AddCoins))
	http.Handle("/coins/transfer", jwt.IsAuthorized(handlers.TransferCoins))
	http.Handle("/coins/balance", jwt.IsAuthorized(handlers.GetCoinBalance))
}
