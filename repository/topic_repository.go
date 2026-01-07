package repository

import (
	"errors"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/db"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/models"
)

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

func GetTopicOwner(topicID int) (int, error) {
	var ownerID int
	err := db.DB.QueryRow(
		"SELECT user_id FROM topics WHERE id = ?",
		topicID,
	).Scan(&ownerID)

	if err != nil {
		return 0, errors.New("topic not found")
	}
	return ownerID, nil
}

func GetTopicByID(topicID int) (models.Topic, error) {
	var t models.Topic
	err := db.DB.QueryRow(
		"SELECT id, title, user_id FROM topics WHERE id = ?",
		topicID,
	).Scan(&t.ID, &t.Title, &t.UserID)

	if err != nil {
		return t, errors.New("topic not found")
	}

	return t, nil
}

func UpdateTopic(topicID int, title string) error {
	_, err := db.DB.Exec(
		"UPDATE topics SET title = ? WHERE id = ?",
		title, topicID,
	)
	return err
}

func DeleteTopic(id, userID int) error {
	ownerID, err := GetTopicOwner(id)
	if err != nil {
		return err
	}

	if ownerID != userID {
		return errors.New("forbidden")
	}

	_, err = db.DB.Exec("DELETE FROM topics WHERE id = ?", id)
	return err
}


