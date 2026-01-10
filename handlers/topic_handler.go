package handlers

import (
	"encoding/json"
	"net/http"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/repository"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/middleware"
)

func CreateTopic(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var req struct {
		Title string `json:"title"`
	}
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

func GetAllTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := repository.GetAllTopics()
	if err != nil {
		http.Error(w, "Failed to fetch topics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}

func GetTopic(w http.ResponseWriter, r *http.Request, topicID int) {
	topic, err := repository.GetTopicByID(topicID)
	if err != nil {
		http.Error(w, "Topic not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topic)
}

func UpdateTopic(w http.ResponseWriter, r *http.Request, topicID int) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ownerID, err := repository.GetTopicOwner(topicID)
	if err != nil {
		http.Error(w, "Topic not found", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := repository.UpdateTopic(topicID, req.Title); err != nil {
		http.Error(w, "Failed to update topic", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteTopic(w http.ResponseWriter, r *http.Request, topicID int) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	err := repository.DeleteTopic(topicID, userID)
	if err != nil {
		if err.Error() == "forbidden" {
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			http.Error(w, "Topic not found", http.StatusNotFound)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
