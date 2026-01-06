package models

import "time"

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	TopicID   int       `json:"topic_id"`
	UserID    int       `json:"user_id"`
	Likes     int       `json:"likes"`
	CreatedAt time.Time `json:"created_at"`
}
