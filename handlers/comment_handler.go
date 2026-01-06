package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/middleware"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/repository"
)

type CreateCommentRequest struct {
	Content string `json:"content"`
	PostID  int    `json:"post_id"`
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	commentID, err := repository.CreateComment(req.Content, req.PostID, userID)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{
		"comment_id": commentID,
	})
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	commentID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	ownerID, err := repository.GetCommentOwner(commentID)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := repository.UpdateComment(commentID, req.Content); err != nil {
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	commentID, _ := strconv.Atoi(r.URL.Query().Get("id"))

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

func LikeComment(w http.ResponseWriter, r *http.Request) {
	commentID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	if err := repository.LikeComment(commentID); err != nil {
		http.Error(w, "Failed to like comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UnlikeComment(w http.ResponseWriter, r *http.Request) {
	commentID, _ := strconv.Atoi(r.URL.Query().Get("id")) 
	if err := repository.UnlikeComment(commentID); err != nil {
		http.Error(w, "Failed to unlike comment", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetCommentsByPost(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	comments, err := repository.GetCommentsByPost(postID)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}