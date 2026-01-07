package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/middleware"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/repository"
)

type CreateTopicRequest struct {
	Title string `json:"title"`
}

type UpdateTopicRequest struct {
	Title string `json:"title"`
}

func CreateTopic(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

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

func GetAllTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := repository.GetAllTopics()
	if err != nil {
		http.Error(w, "Failed to fetch topics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}

func GetTopic(w http.ResponseWriter, r *http.Request) {
	topicID, err := extractIDFromPath(r)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	topic, err := repository.GetTopicByID(topicID)
	if err != nil {
		http.Error(w, "Topic not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topic)
}

func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	topicID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	ownerID, err := repository.GetTopicOwner(topicID)
	if err != nil {
		http.Error(w, "Topic not found", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req UpdateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = repository.UpdateTopic(topicID, req.Title)
	if err != nil {
		http.Error(w, "Failed to update topic", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	topicID, _ := strconv.Atoi(r.URL.Query().Get("id"))

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
