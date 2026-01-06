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
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	res, err := db.DB.Exec(
		"INSERT INTO comments (content, post_id, user_id) VALUES (?, ?, ?)",
		comment.Content, comment.PostID, userID,
	)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	commentID, _ := res.LastInsertId() 

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{
		"comment_id": int(commentID),
	})
}
