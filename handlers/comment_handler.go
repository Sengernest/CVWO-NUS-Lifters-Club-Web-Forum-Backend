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
}

type UpdateCommentRequest struct {
    Content string `json:"content"`
}

// CreateComment handles POST /posts/{postID}/comments
func CreateComment(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.UserIDKey).(int)
    postID := r.Context().Value("postID").(int) // injected from main.go

    var req CreateCommentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    commentID, err := repository.CreateComment(req.Content, postID, userID)
    if err != nil {
        http.Error(w, "Failed to create comment", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]int{
        "comment_id": commentID,
    })
}

// GetCommentsByPost handles GET /posts/{postID}/comments
func GetCommentsByPost(w http.ResponseWriter, r *http.Request) {
    postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    comments, err := repository.GetCommentsByPost(postID)
    if err != nil {
        http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(comments)
}

// UpdateComment handles PUT /comments/{id}
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

// DeleteComment handles DELETE /comments/{id}
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

func ToggleLikeComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	commentID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	liked, err := repository.ToggleCommentLike(commentID, userID)
	if err != nil {
		http.Error(w, "Failed to toggle comment like", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{
		"liked": liked,
	})
}
