package handlers

import (
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/db"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/middleware"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// EditPost allows the owner to edit their post
func EditPost(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	// Ensure user owns the post
	var ownerID int
	err = db.DB.QueryRow("SELECT user_id FROM posts WHERE id = ?", postID).Scan(&ownerID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "You can only edit your own posts", http.StatusUnauthorized)
		return
	}

	var post struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	json.NewDecoder(r.Body).Decode(&post)

	_, err = db.DB.Exec(
		"UPDATE posts SET title = ?, content = ? WHERE id = ?",
		post.Title, post.Content, postID,
	)

	if err != nil {
		http.Error(w, "Failed to edit post", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Post updated successfully")
}

// EditComment allows the owner to edit their comment
func EditComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	commentIDStr := r.URL.Query().Get("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment id", http.StatusBadRequest)
		return
	}

	// Ensure user owns the comment
	var ownerID int
	err = db.DB.QueryRow("SELECT user_id FROM comments WHERE id = ?", commentID).Scan(&ownerID)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "You can only edit your own comments", http.StatusUnauthorized)
		return
	}

	var comment struct {
		Content string `json:"content"`
	}

	json.NewDecoder(r.Body).Decode(&comment)

	_, err = db.DB.Exec(
		"UPDATE comments SET content = ? WHERE id = ?",
		comment.Content, commentID,
	)

	if err != nil {
		http.Error(w, "Failed to edit comment", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Comment updated successfully")
}
