package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/middleware"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/repository"
)

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	TopicID int    `json:"topic_id"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	postID, err := repository.CreatePost(
		req.Title,
		req.Content,
		req.TopicID,
		userID,
	)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{
		"post_id": postID,
	})
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

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	postID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	ownerID, err := repository.GetPostOwner(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = repository.UpdatePost(postID, req.Title, req.Content)
	if err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	postID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	ownerID, err := repository.GetPostOwner(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err = repository.DeletePost(postID)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	err := repository.LikePost(postID)
	if err != nil {
		http.Error(w, "Failed to like post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
