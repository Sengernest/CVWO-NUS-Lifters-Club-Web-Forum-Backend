package handlers

import (
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/db"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/middleware"
	"encoding/json"
	"net/http"
)

type Comment struct {
	Content string `json:"content"`
	PostID  int    `json:"post_id"`
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var comment Comment
	json.NewDecoder(r.Body).Decode(&comment)

	_, err := db.DB.Exec(
		"INSERT INTO comments (content, post_id, user_id) VALUES (?, ?, ?)",
		comment.Content, comment.PostID, userID,
	)

	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
