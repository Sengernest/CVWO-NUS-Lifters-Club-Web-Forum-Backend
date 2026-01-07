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
		`SELECT p.id, p.title, p.content, p.topic_id, p.user_id, u.username, p.likes, p.created_at
		 FROM posts p
		 JOIN users u ON p.user_id = u.id
		 ORDER BY p.id DESC`,
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
			&p.Username, 
			&p.Likes,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}


func GetPostsByTopic(topicID int) ([]models.Post, error) {
	rows, err := db.DB.Query(
		`SELECT p.id, p.title, p.content, p.topic_id, p.user_id, u.username, p.likes, p.created_at
		 FROM posts p
		 JOIN users u ON p.user_id = u.id
		 WHERE p.topic_id = ?
		 ORDER BY p.created_at DESC`,
		topicID,
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
			&p.Username,
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

func TogglePostLike(postID, userID int) (bool, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return false, err
	}

	var exists int
	err = tx.QueryRow(
		`SELECT 1 FROM post_likes WHERE post_id = ? AND user_id = ?`,
		postID, userID,
	).Scan(&exists)

	if err == nil {
	
		_, err = tx.Exec(
			`DELETE FROM post_likes WHERE post_id = ? AND user_id = ?`,
			postID, userID,
		)
		if err != nil {
			tx.Rollback()
			return false, err
		}

		_, err = tx.Exec(
			`UPDATE posts SET likes = likes - 1 WHERE id = ? AND likes > 0`,
			postID,
		)
		if err != nil {
			tx.Rollback()
			return false, err
		}

		tx.Commit()
		return false, nil 
	}

	_, err = tx.Exec(
		`INSERT INTO post_likes (post_id, user_id) VALUES (?, ?)`,
		postID, userID,
	)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	_, err = tx.Exec(
		`UPDATE posts SET likes = likes + 1 WHERE id = ?`,
		postID,
	)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil 
}
