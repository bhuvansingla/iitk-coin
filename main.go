package main

import (
	"fmt"
	"net/http"
)

func foo(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Works!")
}

func main() {
	fmt.Println("Hello World")

	http.HandleFunc("/", foo)

	http.ListenAndServe(":8080", nil)
}
