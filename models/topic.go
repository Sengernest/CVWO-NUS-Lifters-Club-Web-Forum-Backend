package models

type Topic struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	UserID int    `json:"user_id"`
}
