package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func ConnectDatabase() {
	var err error

	DB, err = sql.Open("sqlite", "forum.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	createTables()
}

func createTables() {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	topicTable := `
	CREATE TABLE IF NOT EXISTS topics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	postTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		topic_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		likes INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (topic_id) REFERENCES topics(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	commentTable := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		likes INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	execTable(userTable)
	execTable(topicTable)
	execTable(postTable)
	execTable(commentTable)
}

func execTable(query string) {
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}
