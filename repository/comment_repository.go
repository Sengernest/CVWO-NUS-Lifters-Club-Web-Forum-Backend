package repository

import (
	"errors"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/db"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/models"
)

func CreateComment(content string, postID, userID int) (int, error) {
	res, err := db.DB.Exec(
		"INSERT INTO comments (content, post_id, user_id) VALUES (?, ?, ?)",
		content, postID, userID,
	)
	if err != nil {
		return 0, err
	}

	commentID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(commentID), nil
}

func GetCommentOwner(commentID int) (int, error) {
	var ownerID int
	err := db.DB.QueryRow(
		"SELECT user_id FROM comments WHERE id = ?",
		commentID,
	).Scan(&ownerID)
	if err != nil {
		return 0, errors.New("comment not found")
	}
	return ownerID, nil
}

func GetCommentsByPost(postID int) ([]models.Comment, error) {
	rows, err := db.DB.Query(
		"SELECT id, content, post_id, user_id, likes, created_at FROM comments WHERE post_id = ? ORDER BY id ASC",
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.Content, &c.PostID, &c.UserID, &c.Likes, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}


func UpdateComment(commentID int, content string) error {
	_, err := db.DB.Exec(
		"UPDATE comments SET content = ? WHERE id = ?",
		content, commentID,
	)
	return err
}

func DeleteComment(commentID int) error {
	_, err := db.DB.Exec(
		"DELETE FROM comments WHERE id = ?",
		commentID,
	)
	return err
}

func LikeComment(commentID int) error {
	_, err := db.DB.Exec(
		"UPDATE comments SET likes = likes + 1 WHERE id = ?",
		commentID,
	)
	return err
}
func UnlikeComment(commentID int) error {
	_, err := db.DB.Exec(
		"UPDATE comments SET likes = likes - 1 WHERE id = ? AND likes > 0",
		commentID,
	)
	return err
}				
