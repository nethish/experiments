package main

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Response{Message: "Hello, World!"})
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Response{Message: "This is a sample Go server"})
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		json.NewEncoder(w).Encode(Response{Message: "Logged in successfully"})
	})

	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		json.NewEncoder(w).Encode(Response{Message: "dashboard"})
	})

	http.ListenAndServe(":8080", nil)
}
