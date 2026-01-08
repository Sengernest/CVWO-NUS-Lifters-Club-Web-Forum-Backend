package repository

import (
	"errors"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/db"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/models"
)

func CreateComment(content string, postID, userID int) (int, error) {
	res, err := db.DB.Exec("INSERT INTO comments (content, post_id, user_id) VALUES (?, ?, ?)", content, postID, userID)
	if err != nil {
		return 0, err
	}

	id64, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id64), nil
}

func GetCommentOwner(commentID int) (int, error) {
	var ownerID int
	err := db.DB.QueryRow("SELECT user_id FROM comments WHERE id = ?", commentID).Scan(&ownerID)
	if err != nil {
		return 0, errors.New("comment not found")
	}
	return ownerID, nil
}

func GetCommentsByPost(postID int) ([]models.Comment, error) {
	rows, err := db.DB.Query(
		`SELECT c.id, c.content, c.post_id, c.user_id, u.username, c.likes, c.created_at
		 FROM comments c
		 JOIN users u ON c.user_id = u.id
		 WHERE c.post_id = ?
		 ORDER BY c.created_at ASC`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.Content, &c.PostID, &c.UserID, &c.Username, &c.Likes, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func UpdateComment(commentID int, content string) error {
	_, err := db.DB.Exec("UPDATE comments SET content = ? WHERE id = ?", content, commentID)
	return err
}

func DeleteComment(commentID int) error {
	_, err := db.DB.Exec("DELETE FROM comments WHERE id = ?", commentID)
	return err
}

func ToggleCommentLike(commentID, userID int) (bool, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return false, err
	}

	var exists int
	err = tx.QueryRow("SELECT 1 FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&exists)
	if err == nil {
		// Already liked: remove
		if _, err := tx.Exec("DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID); err != nil {
			tx.Rollback()
			return false, err
		}
		if _, err := tx.Exec("UPDATE comments SET likes = likes - 1 WHERE id = ? AND likes > 0", commentID); err != nil {
			tx.Rollback()
			return false, err
		}
		tx.Commit()
		return false, nil
	}

	// Not liked: insert
	if _, err := tx.Exec("INSERT INTO comment_likes (comment_id, user_id) VALUES (?, ?)", commentID, userID); err != nil {
		tx.Rollback()
		return false, err
	}
	if _, err := tx.Exec("UPDATE comments SET likes = likes + 1 WHERE id = ?", commentID); err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}
