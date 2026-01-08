package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/db"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/handlers"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/middleware"
)

func main() {
	// Connect DB
	db.ConnectDatabase()

	mux := http.NewServeMux()

	// --- Public Auth Routes ---
	mux.HandleFunc("/register", middleware.Cors(handlers.Register))
	mux.HandleFunc("/login", middleware.Cors(handlers.Login))

	// --- Topics ---
	mux.HandleFunc("/topics", middleware.Cors(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetAllTopics(w, r)
		case http.MethodPost:
			middleware.AuthMiddleware(handlers.CreateTopic)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/topics/", middleware.Cors(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 2 {
			http.Error(w, "Invalid topics route", http.StatusBadRequest)
			return
		}

		topicID, err := strconv.Atoi(parts[1])
		if err != nil {
			http.Error(w, "Invalid topic ID", http.StatusBadRequest)
			return
		}

		// /topics/:id/posts
		if len(parts) == 3 && parts[2] == "posts" {
			switch r.Method {
			case http.MethodGet:
				handlers.GetPostsByTopic(w, r, topicID)
			case http.MethodPost:
				middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
					ctx := context.WithValue(r.Context(), "topicID", topicID)
					r = r.WithContext(ctx)
					handlers.CreatePost(w, r)
				})(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// /topics/:id
		switch r.Method {
		case http.MethodGet:
			handlers.GetTopic(w, r, topicID)
		case http.MethodPut:
			middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.UpdateTopic(w, r, topicID)
			})(w, r)
		case http.MethodDelete:
			middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.DeleteTopic(w, r, topicID)
			})(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// --- Posts ---
	mux.HandleFunc("/posts", middleware.Cors(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetAllPosts(w, r)
		case http.MethodPost:
			middleware.AuthMiddleware(handlers.CreatePost)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/posts/", middleware.Cors(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 2 {
			http.Error(w, "Invalid posts route", http.StatusBadRequest)
			return
		}

		postID, err := strconv.Atoi(parts[1])
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// /posts/:id/like
		if len(parts) == 3 && parts[2] == "like" && r.Method == http.MethodPost {
			middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.ToggleLikePost(w, r, postID)
			})(w, r)
			return
		}

		// /posts/:id/comments
		if len(parts) == 3 && parts[2] == "comments" {
			switch r.Method {
			case http.MethodGet:
				handlers.GetCommentsByPost(w, r, postID)
			case http.MethodPost:
				middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
					ctx := context.WithValue(r.Context(), "postID", postID)
					r = r.WithContext(ctx)
					handlers.CreateComment(w, r)
				})(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// /posts/:id
		switch r.Method {
		case http.MethodPut:
			middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.UpdatePost(w, r, postID)
			})(w, r)
		case http.MethodDelete:
			middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.DeletePost(w, r, postID)
			})(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// --- Comments ---
	mux.HandleFunc("/comments/", middleware.Cors(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 2 {
			http.Error(w, "Invalid comment route", http.StatusBadRequest)
			return
		}

		commentID, err := strconv.Atoi(parts[1])
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		// /comments/:id/like
		if len(parts) == 3 && parts[2] == "like" && r.Method == http.MethodPost {
			middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.ToggleLikeComment(w, r, commentID)
			})(w, r)
			return
		}

		// /comments/:id
		switch r.Method {
		case http.MethodPut:
			middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.UpdateComment(w, r, commentID)
			})(w, r)
		case http.MethodDelete:
			middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				handlers.DeleteComment(w, r, commentID)
			})(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// --- Default ---
	mux.HandleFunc("/", middleware.Cors(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "NUS Lifters Club backend running with SQLite")
	}))

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
