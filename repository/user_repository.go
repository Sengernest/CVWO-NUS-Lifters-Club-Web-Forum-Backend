package repository

import (
	"database/sql"
	"errors"

	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/db"
	"CVWO-NUS-Lifters-Club-Web-Forum-Backend/backend/models"
)

func CreateUser(username, passwordHash string) error {
	_, err := db.DB.Exec(
		"INSERT INTO users (username, password_hash) VALUES (?, ?)",
		username,
		passwordHash,
	)
	return err
}

func GetUserByUsername(username string) (*models.User, string, error) {
	var user models.User
	var passwordHash string

	err := db.DB.QueryRow(
		"SELECT id, username, password_hash, created_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &passwordHash, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, "", errors.New("user not found")
	}

	return &user, passwordHash, err
}
