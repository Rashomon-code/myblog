package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestInitAPP_success(t *testing.T) {
	if testDB == nil {
		t.Fatalf("testDB is nil")
	}

	err := testDB.Ping()
	if err != nil {
		t.Fatalf("接続に失敗しました: %v", err)
	}
}

func TestInitSQL_CreatesTables(t *testing.T) {
	tables := []string{"users", "posts", "user_profiles"}
	query := `SELECT EXISTS (
		SELECT FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_name = $1
	)`
	for _, table := range tables {
		t.Run(table, func(t *testing.T) {
			var exists bool
			row := testDB.QueryRow(query, table)
			err := row.Scan(&exists)
			if err != nil {
				t.Fatalf("テーブル %s が確認できませんでした: %v", table, err)
			}
			assert.True(t, exists)
		})
	}
}

func TestInitAdmin_CreatesAdmin(t *testing.T) {
	var username string
	var role string

	err := testDB.QueryRow(`
		SELECT username, role
		FROM users
		WHERE role = 'admin'
		LIMIT 1
	`).Scan(&username, &role)
	if err != nil {
		t.Fatalf("検索できませんでした: %v", err)
	}

	assert.Equal(t, "admin", username)
	assert.Equal(t, "admin", role)
}

func TestInitAdmin_HasAdmin(t *testing.T) {
	var before int
	var after int

	err := testDB.QueryRow(`
		SELECT COUNT(*)
		FROM users
		WHERE role = 'admin'
	`).Scan(&before)
	if err != nil {
		t.Fatalf("検索できませんでした: %v", err)
	}

	if before != 1 {
		t.Fatalf("expected admin count %q, got %q", 1, before)
	}

	err = initAdmin(testDB)
	if err != nil {
		t.Fatalf("admin 初期化できませんでした: %v", err)
	}

	err = testDB.QueryRow(`
		SELECT COUNT(*)
		FROM users
		WHERE role = 'admin'
	`).Scan(&after)
	if err != nil {
		t.Fatalf("検索できませんでした: %v", err)
	}

	if after != 1 {
		t.Fatalf("expected admin count %q, got %q", 1, after)
	}
}

func TestInitAdmin_PasswordISCorrect(t *testing.T) {
	var passwordHash string

	err := testDB.QueryRow(`
		SELECT password_hash	
		FROM users
		WHERE username = 'admin'
	`).Scan(&passwordHash)
	if err != nil {
		t.Fatalf("検索できませんでした: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte("admin1234"))
	if err != nil {
		t.Fatalf("パスワードが正しくありません :%v", err)
	}
}
