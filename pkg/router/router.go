package router

import (
	"net/http"

	"github.com/bhuvansingla/iitk-coin/pkg/handlers"
)

func SetRoutes() {
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/signup", handlers.Signup)
	http.HandleFunc("/coins/add", handlers.AddCoins)
	http.HandleFunc("/coins/transfer", handlers.TransferCoins)
	http.HandleFunc("/coins/balance", handlers.GetCoinBalance)
	// http.Handle("/access", jwt.IsAuthorized(secretPage))
}
