package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/db"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/handlers"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/middleware"
)

// extractIDFromPath parses the last segment as ID
func extractIDFromPath(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	idStr := parts[len(parts)-1]
	return strconv.Atoi(idStr)
}

func main() {
	// Connect DB
	db.ConnectDatabase()

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/register", middleware.Cors(handlers.Register))
	mux.HandleFunc("/login", middleware.Cors(handlers.Login))

	// Posts collection
	mux.HandleFunc("/posts", middleware.Cors(middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetAllPosts(w, r)
		case http.MethodPost:
			handlers.CreatePost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Topics collection (GET / POST)
mux.HandleFunc("/topics", middleware.Cors(middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handlers.GetAllTopics(w, r)
	case http.MethodPost:
		handlers.CreateTopic(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
})))

// Topics item route (DELETE /topics/{id})
mux.HandleFunc("/topics/", middleware.Cors(middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path: /topics/{id}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Missing topic ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	// Put ID into query param so DeleteTopic handler can use it
	q := r.URL.Query()
	q.Set("id", strconv.Itoa(id))
	r.URL.RawQuery = q.Encode()

	handlers.DeleteTopic(w, r)
})))



	// Posts item routes + comments
	mux.HandleFunc("/posts/", middleware.Cors(middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// If /posts/{id}/like
		if strings.HasSuffix(path, "/like") && r.Method == http.MethodPost {
			id, err := extractIDFromPath(r)
			if err != nil {
				http.Error(w, "Invalid post ID", http.StatusBadRequest)
				return
			}
			r.URL.Query().Set("id", strconv.Itoa(id))
			handlers.LikePost(w, r)
			return
		}

		// If /posts/{id}/comments
		if strings.HasSuffix(path, "/comments") {
			id, err := extractIDFromPath(r)
			if err != nil {
				http.Error(w, "Invalid post ID", http.StatusBadRequest)
				return
			}
			r.URL.Query().Set("id", strconv.Itoa(id))

			switch r.Method {
			case http.MethodPost:
				handlers.CreateComment(w, r)
			case http.MethodGet:
				handlers.GetCommentsByPost(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Otherwise, /posts/{id} → Update/Delete
		id, err := extractIDFromPath(r)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		r.URL.Query().Set("id", strconv.Itoa(id))

		switch r.Method {
		case http.MethodPut:
			handlers.UpdatePost(w, r)
		case http.MethodDelete:
			handlers.DeletePost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Comments item routes
	mux.HandleFunc("/comments/", middleware.Cors(middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		id, err := extractIDFromPath(r)
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}
		r.URL.Query().Set("id", strconv.Itoa(id))

		if strings.HasSuffix(path, "/like") && r.Method == http.MethodPost {
			handlers.LikeComment(w, r)
			return
		}

		switch r.Method {
		case http.MethodPut:
			handlers.UpdateComment(w, r)
		case http.MethodDelete:
			handlers.DeleteComment(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Root route
	mux.HandleFunc("/", middleware.Cors(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "NUS Lifters Club backend running with SQLite")
	}))

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Server failed:", err)
	}
}
