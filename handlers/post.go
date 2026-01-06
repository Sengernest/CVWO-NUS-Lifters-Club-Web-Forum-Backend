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
	json.NewDecoder(r.Body).Decode(&post)

	_, err := db.DB.Exec(
		"INSERT INTO posts (title, content, topic_id, user_id) VALUES (?, ?, ?, ?)",
		post.Title, post.Content, post.TopicID, userID,
	)

	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
