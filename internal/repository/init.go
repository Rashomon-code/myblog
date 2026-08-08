package repository

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitSQL() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("データベースを開く際にエラーが発生しました: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("データベースの接続に問題が起きました: %w", err)
	}
	fmt.Println("PostgreSQL に接続済み")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS posts(
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`
	//FOREIGN KEY, PRIMARY KEY など table constraints は最後に書かなければなりません。

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("テーブルが作成できませんでした: %w", err)
	}
	fmt.Println("テーブルを作成済み。")

	return db, nil
}
