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

    // Public routes
    http.HandleFunc("/register", middleware.Cors(handlers.Register))
    http.HandleFunc("/login", middleware.Cors(handlers.Login))

    // Protected routes (Auth + CORS)
    http.HandleFunc("/posts", middleware.Cors(middleware.AuthMiddleware(handlers.CreatePost)))
	http.HandleFunc("/getposts", middleware.Cors(middleware.AuthMiddleware(handlers.GetAllPosts)))
    http.HandleFunc("/comments", middleware.Cors(middleware.AuthMiddleware(handlers.CreateComment)))
    http.HandleFunc("/delete-post", middleware.Cors(middleware.AuthMiddleware(handlers.DeletePost)))
    http.HandleFunc("/delete-comment", middleware.Cors(middleware.AuthMiddleware(handlers.DeleteComment)))
    http.HandleFunc("/edit-post", middleware.Cors(middleware.AuthMiddleware(handlers.EditPost)))
    http.HandleFunc("/edit-comment", middleware.Cors(middleware.AuthMiddleware(handlers.EditComment)))
    http.HandleFunc("/like-post", middleware.Cors(middleware.AuthMiddleware(handlers.LikePost)))
    http.HandleFunc("/like-comment", middleware.Cors(middleware.AuthMiddleware(handlers.LikeComment)))

    // Root route
    http.HandleFunc("/", middleware.Cors(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "NUS Lifters Club backend running with SQLite")
    }))

    fmt.Println("Server running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
