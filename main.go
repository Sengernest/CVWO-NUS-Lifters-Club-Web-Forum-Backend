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

// extractIDFromPath parses the last segment as ID
func extractIDFromPath(r *http.Request) (int, error) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	return strconv.Atoi(parts[len(parts)-1])
}

func main() {
	// Connect DB
	db.ConnectDatabase()

	mux := http.NewServeMux()

	// ================= AUTH =================
	mux.HandleFunc("/register", middleware.Cors(handlers.Register))
	mux.HandleFunc("/login", middleware.Cors(handlers.Login))

	// ================= TOPICS COLLECTION =================
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

	// ================= TOPICS + POSTS UNDER TOPIC =================
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

		// ===== /topics/{id}/posts =====
		if len(parts) == 3 && parts[2] == "posts" {
			q := r.URL.Query()
			q.Set("topic_id", strconv.Itoa(topicID))
			r.URL.RawQuery = q.Encode()

			switch r.Method {
			case http.MethodGet:
				handlers.GetPostsByTopic(w, r)
			case http.MethodPost:
				middleware.AuthMiddleware(handlers.CreatePost)(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		  // /topics/{id} 
		if len(parts) == 2 && r.Method == http.MethodGet {
			handlers.GetTopic(w, r) 
			return
		}

			// ===== /topics/{id} =====
		q := r.URL.Query()
		q.Set("id", strconv.Itoa(topicID))
		r.URL.RawQuery = q.Encode()

		switch r.Method {
		case http.MethodPut:
			middleware.AuthMiddleware(handlers.UpdateTopic)(w, r)
		case http.MethodDelete:
			middleware.AuthMiddleware(handlers.DeleteTopic)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// ================= POSTS COLLECTION =================
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

	// ================= POSTS ITEM + COMMENTS =================
	mux.HandleFunc("/posts/", middleware.Cors(func(w http.ResponseWriter, r *http.Request) {
    path := strings.Trim(r.URL.Path, "/")
    parts := strings.Split(path, "/")

    if len(parts) < 2 || parts[0] != "posts" {
        http.Error(w, "Invalid posts route", http.StatusBadRequest)
        return
    }

    postID, err := strconv.Atoi(parts[1])
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }
   q := r.URL.Query()
q.Set("id", strconv.Itoa(postID))
r.URL.RawQuery = q.Encode() // this actually updates the query for handlers


    // /posts/{id}/like (needs auth)
    if len(parts) == 3 && parts[2] == "like" && r.Method == http.MethodPost {
        middleware.AuthMiddleware(handlers.ToggleLikePost)(w, r)
        return
    }

// /posts/{id}/comments
if len(parts) == 3 && parts[2] == "comments" {
    postID := postID // already parsed above

    // Inject postID into query for GetCommentsByPost
    q := r.URL.Query()
    q.Set("post_id", strconv.Itoa(postID))
    r.URL.RawQuery = q.Encode()

    switch r.Method {
    case http.MethodGet:
        handlers.GetCommentsByPost(w, r) // PUBLIC
    case http.MethodPost:
        // inject postID into context for CreateComment
        ctx := r.Context()
        ctx = context.WithValue(ctx, "postID", postID)
        r = r.WithContext(ctx)

        middleware.AuthMiddleware(handlers.CreateComment)(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
    return
}

    // /posts/{id} (edit/delete requires auth)
    switch r.Method {
    case http.MethodPut:
        middleware.AuthMiddleware(handlers.UpdatePost)(w, r)
    case http.MethodDelete:
        middleware.AuthMiddleware(handlers.DeletePost)(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}))

// ================= COMMENTS ITEM =================
mux.HandleFunc("/comments/", middleware.Cors(func(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    if len(parts) < 2 || parts[0] != "comments" {
        http.Error(w, "Invalid comment route", http.StatusBadRequest)
        return
    }

    commentID, err := strconv.Atoi(parts[1])
    if err != nil {
        http.Error(w, "Invalid comment ID", http.StatusBadRequest)
        return
    }

    // Inject comment ID into query for handlers
    q := r.URL.Query()
    q.Set("id", strconv.Itoa(commentID))
    r.URL.RawQuery = q.Encode()

    // Handle /comments/{id}/like
    if len(parts) == 3 && parts[2] == "like" && r.Method == http.MethodPost {
        middleware.AuthMiddleware(handlers.ToggleLikeComment)(w, r)
        return
    }

    // CRUD for comments
    switch r.Method {
    case http.MethodPut:
        middleware.AuthMiddleware(handlers.UpdateComment)(w, r)
    case http.MethodDelete:
        middleware.AuthMiddleware(handlers.DeleteComment)(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}))


	// ================= ROOT =================
	mux.HandleFunc("/", middleware.Cors(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "NUS Lifters Club backend running with SQLite")
	}))

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
