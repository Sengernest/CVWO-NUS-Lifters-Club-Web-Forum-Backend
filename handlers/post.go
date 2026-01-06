package handlers

import (
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/db"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/middleware"
	"encoding/json"
	"net/http"
)

type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	TopicID int    `json:"topic_id"`
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	res, err := db.DB.Exec(
		"INSERT INTO posts (title, content, topic_id, user_id) VALUES (?, ?, ?, ?)",
		post.Title, post.Content, post.TopicID, userID,
	)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	postID, _ := res.LastInsertId() 

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{
		"post_id": int(postID),
	})
}

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	// Fetch posts from DB
	rows, err := db.DB.Query("SELECT id, title, content, topic_id FROM posts ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type PostResponse struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		TopicID int    `json:"topic_id"`
	}

	var posts []PostResponse
	for rows.Next() {
		var p PostResponse
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.TopicID); err != nil {
			http.Error(w, "Error scanning posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}



