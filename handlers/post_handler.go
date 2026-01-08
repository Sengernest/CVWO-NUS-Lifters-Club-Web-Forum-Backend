package handlers

import (
	"encoding/json"
	"net/http"

	
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/repository"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/middleware"
)

// CreatePost expects AuthMiddleware
func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	topicID, _ := r.Context().Value("topicID").(int) // optional, set via main

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" || req.Content == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	postID, err := repository.CreatePost(req.Title, req.Content, topicID, userID)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"post_id": postID})
}

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := repository.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetPostsByTopic(w http.ResponseWriter, r *http.Request, topicID int) {
	posts, err := repository.GetPostsByTopic(topicID)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func UpdatePost(w http.ResponseWriter, r *http.Request, postID int) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" || req.Content == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ownerID, err := repository.GetPostOwner(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	if ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := repository.UpdatePost(postID, req.Title, req.Content); err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeletePost(w http.ResponseWriter, r *http.Request, postID int) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	ownerID, err := repository.GetPostOwner(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	if ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := repository.DeletePost(postID); err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ToggleLikePost(w http.ResponseWriter, r *http.Request, postID int) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	liked, err := repository.TogglePostLike(postID, userID)
	if err != nil {
		http.Error(w, "Failed to toggle like", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"liked": liked})
}
