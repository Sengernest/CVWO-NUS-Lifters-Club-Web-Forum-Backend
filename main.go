package main

import (
	"fmt"
	"net/http"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/db"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/handlers"
)

func main() {
	db.ConnectDatabase()
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "NUS Lifters Club backend running with SQLite")
	})

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
