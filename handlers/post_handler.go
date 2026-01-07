package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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


// extractIDFromPath parses the last segment of the URL path as an integer ID
func extractIDFromPath(r *http.Request) (int, error) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid path")
	}
	return strconv.Atoi(parts[len(parts)-1])
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

func GetPostsByTopic(w http.ResponseWriter, r *http.Request) {
	topicID, err := strconv.Atoi(r.URL.Query().Get("topic_id"))
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	posts, err := repository.GetPostsByTopic(topicID)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}


// UpdatePost handles PUT /posts/{id}
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Extract post ID from path
	postID, err := extractIDFromPath(r)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	fmt.Println("UpdatePost called for postID:", postID, "by userID:", userID)

	ownerID, err := repository.GetPostOwner(postID)
	if err != nil {
		fmt.Println("GetPostOwner failed:", err)
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

	if err := repository.UpdatePost(postID, req.Title, req.Content); err != nil {
		fmt.Println("UpdatePost error:", err)
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	fmt.Println("Post updated successfully:", postID)
	w.WriteHeader(http.StatusOK)
}

// DeletePost handles DELETE /posts/{id}
func DeletePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Extract post ID from path
	postID, err := extractIDFromPath(r)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	fmt.Println("DeletePost called for postID:", postID, "by userID:", userID)

	ownerID, err := repository.GetPostOwner(postID)
	if err != nil {
		fmt.Println("GetPostOwner failed:", err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := repository.DeletePost(postID); err != nil {
		fmt.Println("DeletePost error:", err)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	fmt.Println("Post deleted successfully:", postID)
	w.WriteHeader(http.StatusNoContent)
}

func ToggleLikePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	postID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	liked, err := repository.TogglePostLike(postID, userID)
	if err != nil {
		http.Error(w, "Failed to toggle like", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{
		"liked": liked,
	})
}


