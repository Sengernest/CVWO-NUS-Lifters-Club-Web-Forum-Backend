package models

import "time"

type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	Likes     int       `json:"likes"`
	CreatedAt time.Time `json:"created_at"`
}
