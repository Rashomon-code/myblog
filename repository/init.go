package repository

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func InitSQL() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./blog.db")
	if err != nil {
		return nil, fmt.Errorf("打開數據庫失敗: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("無法連接到數據庫: %w", err)
	}
	fmt.Println("已連接 SQLite")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS posts(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("創建表失敗: %w", err)
	}
	fmt.Println("表已創建")

	return db, nil
}
