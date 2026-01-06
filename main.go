package main

import (
	"fmt"
	"net/http"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/db"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/handlers"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/middleware" 
	
)

func main() {
	db.ConnectDatabase()
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)

	// Protected routes
	http.HandleFunc("/posts", middleware.AuthMiddleware(handlers.CreatePost))
	http.HandleFunc("/comments", middleware.AuthMiddleware(handlers.CreateComment))
	http.HandleFunc("/delete-post", middleware.AuthMiddleware(handlers.DeletePost))
	http.HandleFunc("/delete-comment", middleware.AuthMiddleware(handlers.DeleteComment))


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "NUS Lifters Club backend running with SQLite")
	})

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
