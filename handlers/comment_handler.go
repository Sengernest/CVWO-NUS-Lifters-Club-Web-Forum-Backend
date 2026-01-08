package handlers

import (
	"encoding/json"
	"net/http"


	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/repository"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/middleware"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	postID := r.Context().Value("postID").(int)

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	commentID, err := repository.CreateComment(req.Content, postID, userID)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"comment_id": commentID})
}

func GetCommentsByPost(w http.ResponseWriter, r *http.Request, postID int) {
	comments, err := repository.GetCommentsByPost(postID)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func UpdateComment(w http.ResponseWriter, r *http.Request, commentID int) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ownerID, err := repository.GetCommentOwner(commentID)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}
	if ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := repository.UpdateComment(commentID, req.Content); err != nil {
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteComment(w http.ResponseWriter, r *http.Request, commentID int) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	ownerID, err := repository.GetCommentOwner(commentID)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}
	if ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := repository.DeleteComment(commentID); err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ToggleLikeComment(w http.ResponseWriter, r *http.Request, commentID int) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	liked, err := repository.ToggleCommentLike(commentID, userID)
	if err != nil {
		http.Error(w, "Failed to toggle comment like", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"liked": liked})
}
