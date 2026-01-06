package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/middleware"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/repository"
)

type CreateTopicRequest struct {
	Title string `json:"title"`
}

// CreateTopic creates a new topic owned by the logged-in user
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	ctxUserID := r.Context().Value(middleware.UserIDKey)
	userID, ok := ctxUserID.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	topic, err := repository.CreateTopic(req.Title, userID)
	if err != nil {
		http.Error(w, "Failed to create topic", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topic)
}

// GetAllTopics returns all topics
func GetAllTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := repository.GetAllTopics()
	if err != nil {
		http.Error(w, "Failed to fetch topics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}

// DeleteTopic deletes a topic only if the user owns it
func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	ctxUserID := r.Context().Value(middleware.UserIDKey)
	userID, ok := ctxUserID.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract ID from URL path: /topics/{id}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	err = repository.DeleteTopic(id, userID)
	if err != nil {
		if err.Error() == "forbidden" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		} else if err.Error() == "topic not found" {
			http.Error(w, "Topic not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete topic", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
