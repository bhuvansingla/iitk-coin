package router

import (
	"net/http"

	"github.com/bhuvansingla/iitk-coin/pkg/handlers"
	"github.com/bhuvansingla/iitk-coin/pkg/jwt"
)

func SetRoutes() {
	http.Handle("/", jwt.IsAuthorized(handlers.Index))
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/signup", handlers.Signup)
	http.HandleFunc("/coins/add", handlers.AddCoins)
	http.HandleFunc("/coins/transfer", handlers.TransferCoins)
	http.HandleFunc("/coins/balance", handlers.GetCoinBalance)
	// http.Handle("/access", jwt.IsAuthorized(secretPage))
}
