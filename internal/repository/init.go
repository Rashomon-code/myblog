package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

func InitAPP() (*sql.DB, error) {
	db, err := initSQL()
	if err != nil {
		return nil, err
	}

	err = initAdmin(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initSQL() (*sql.DB, error) {
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
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user'
	);

	CREATE TABLE IF NOT EXISTS posts(
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS user_profiles(
		user_id INTEGER PRIMARY KEY,
		display_name TEXT,
		bio TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
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

func initAdmin(db *sql.DB) error {
	var count int
	selectSQL := `SELECT COUNT(*) FROM users WHERE role = 'admin'`

	row := db.QueryRow(selectSQL)
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		defaultUsername := os.Getenv("INITIAL_ADMIN_USER")
		defaultPassword := os.Getenv("INITIAL_ADMIN_PASS")

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), 4)
		if err != nil {
			return fmt.Errorf("管理者初期化失敗しました: %w", err)
		}

		initAdminSQL := `INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3)`
		_, err = db.Exec(initAdminSQL, defaultUsername, passwordHash, "admin")
		if err != nil {
			return fmt.Errorf("アカウント初期化失敗しました: %w", err)
		}

		log.Println("管理者アカウントが生成されました")
	}

	return nil
}
