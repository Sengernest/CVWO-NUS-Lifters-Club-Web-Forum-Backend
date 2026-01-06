package repository

import (
	"errors"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/db"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/models"
)

// CreateTopic inserts a topic with the owner userID
func CreateTopic(title string, userID int) (models.Topic, error) {
	res, err := db.DB.Exec(
		"INSERT INTO topics (title, user_id) VALUES (?, ?)",
		title, userID,
	)
	if err != nil {
		return models.Topic{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Topic{}, err
	}

	return models.Topic{
		ID:     int(id),
		Title:  title,
		UserID: userID,
	}, nil
}

// GetAllTopics fetches all topics
func GetAllTopics() ([]models.Topic, error) {
	rows, err := db.DB.Query(
		"SELECT id, title, user_id FROM topics ORDER BY id ASC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var t models.Topic
		if err := rows.Scan(&t.ID, &t.Title, &t.UserID); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}

	return topics, nil
}

// DeleteTopic deletes a topic only if the user is the owner
func DeleteTopic(id, userID int) error {
	var ownerID int
	err := db.DB.QueryRow(
		"SELECT user_id FROM topics WHERE id = ?",
		id,
	).Scan(&ownerID)
	if err != nil {
		return errors.New("topic not found")
	}

	if ownerID != userID {
		return errors.New("forbidden")
	}

	_, err = db.DB.Exec("DELETE FROM topics WHERE id = ?", id)
	return err
}
