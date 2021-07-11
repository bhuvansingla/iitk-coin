package server

import "net/http"

func SetCorsPolicy(w *http.ResponseWriter, req *http.Request) {
	print("here")
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

func CORS(endpoint func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		SetCorsPolicy(&w, r)
		if r.Method == "OPTIONS" {
			return
		}
		endpoint(w, r)
	}
}
