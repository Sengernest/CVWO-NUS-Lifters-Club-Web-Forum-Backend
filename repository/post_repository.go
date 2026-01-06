package repository

import (
	"errors"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/db"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/models"
)

func CreatePost(title, content string, topicID, userID int) (int, error) {
	res, err := db.DB.Exec(
		"INSERT INTO posts (title, content, topic_id, user_id) VALUES (?, ?, ?, ?)",
		title, content, topicID, userID,
	)
	if err != nil {
		return 0, err
	}

	postID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

func GetAllPosts() ([]models.Post, error) {
	rows, err := db.DB.Query(
		`SELECT id, title, content, topic_id, user_id, likes, created_at
		 FROM posts
		 ORDER BY id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.TopicID,
			&p.UserID,
			&p.Likes,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func GetPostOwner(postID int) (int, error) {
	var ownerID int
	err := db.DB.QueryRow(
		"SELECT user_id FROM posts WHERE id = ?",
		postID,
	).Scan(&ownerID)

	if err != nil {
		return 0, errors.New("post not found")
	}

	return ownerID, nil
}

func UpdatePost(postID int, title, content string) error {
	_, err := db.DB.Exec(
		"UPDATE posts SET title = ?, content = ? WHERE id = ?",
		title, content, postID,
	)
	return err
}

func DeletePost(postID int) error {
	_, err := db.DB.Exec(
		"DELETE FROM posts WHERE id = ?",
		postID,
	)
	return err
}

func LikePost(postID int) error {
	_, err := db.DB.Exec(
		"UPDATE posts SET likes = likes + 1 WHERE id = ?",
		postID,
	)
	return err
}
func UnlikePost(postID int) error {
	_, err := db.DB.Exec(
		"UPDATE posts SET likes = likes - 1 WHERE id = ? AND likes > 0",
		postID,
	)
	return err
}