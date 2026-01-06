package handlers

import (
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/db"
	"fmt"
	"net/http"
	"strconv"
)

// LikePost increments likes of a post
func LikePost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(
		"UPDATE posts SET likes = likes + 1 WHERE id = ?",
		postID,
	)
	if err != nil {
		http.Error(w, "Failed to like post", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Post liked successfully")
}

// LikeComment increments likes of a comment
func LikeComment(w http.ResponseWriter, r *http.Request) {
	commentIDStr := r.URL.Query().Get("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment id", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(
		"UPDATE comments SET likes = likes + 1 WHERE id = ?",
		commentID,
	)
	if err != nil {
		http.Error(w, "Failed to like comment", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Comment liked successfully")
}
