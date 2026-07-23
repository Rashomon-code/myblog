package repository

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func InitSQL() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./blog.db")
	if err != nil {
		return nil, fmt.Errorf("データベースを開く際にエラーが発生しました: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("データベースの接続に問題が起きました: %w", err)
	}
	fmt.Println("SQLite に接続済み")

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
		return nil, fmt.Errorf("テーブルが作成できませんでした: %w", err)
	}
	fmt.Println("テーブルを作成済み。")

	return db, nil
}
