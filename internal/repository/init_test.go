package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/moby/moby/api/types/network"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
)

func setupTestDB(t *testing.T) func() {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			//初期化して一回再起動されます、一回しかチャックしない場合、 ready の同時に stop されます
			// 13:03:52 🔔 Container is ready: ae49e7e12d95
			// 13:03:52 🐳 Stopping container: ae49e7e12d95
			// wait.ForLog("database system is ready to accept connections").
			// 	WithOccurrence(2).
			// 	WithStartupTimeout(30*time.Second),
			wait.ForSQL("5432/tcp", "pgx", func(host string, port network.Port) string {
				return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", "testuser", "testpass", host, port.Port(), "testdb")
			}).WithStartupTimeout(30*time.Second).WithPollInterval(500*time.Millisecond),
		),
	)
	if err != nil {
		t.Fatalf("Postgres コンテナー起動できませんでした: %v", err)
	}

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")

	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("INITIAL_ADMIN_USER", "admin")
	os.Setenv("INITIAL_ADMIN_PASS", "admin1234")

	return func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("クリアできませんでした: %v", err)
		}
	}
}

func TestInitAPP_success(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	fmt.Println("hello")
	db, err := InitAPP()
	if err != nil {
		t.Fatalf("初期化できませんでした: %v", err)
	}
	defer db.Close()

	tables := []string{"users", "posts", "user_profiles"}
	for _, table := range tables {
		var exists bool
		query := `SELECT EXISTS (
			SELECT FROM information_schema.tables	
			WHERE table_name = $1
		)`

		err := db.QueryRow(query, table).Scan(&exists)
		if err != nil {
			t.Errorf("テーブル %s の検索に失敗しました: %v", table, err)
		}
		if !exists {
			t.Errorf("テーブル %s が存在しませんでした", table)
		}
	}

	var (
		username     string
		passwordHash string
		role         string
	)

	query := `SELECT username, password_hash, role FROM users WHERE role = 'admin'`
	row := db.QueryRow(query)
	err = row.Scan(&username, &passwordHash, &role)
	if err != nil {
		t.Fatalf("検索できませんでした: %v", err)
	}
	if username != "admin" {
		t.Errorf("expected username %q, got %q", "admin", username)
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte("admin1234"))
	if err != nil {
		t.Errorf("パスワード検証できませんでした: %v", err)
	}
}

func TestInitAPP_Idempotency(t *testing.T) {
	setupTestDB(t)

	db, err := InitAPP()
	if err != nil {
		t.Fatalf("初期化できませんでした: %v", err)
	}

	_, err = InitAPP()
	if err != nil {
		t.Fatalf("二度目の初期化に失敗しました: %v", err)
	}

	var count int
	query := `SELECT COUNT(*) FROM users WHERE role = 'admin'`
	row := db.QueryRow(query)
	err = row.Scan(&count)
	if err != nil || count != 1 {
		t.Errorf("expect no err and count %d, got %v", 1, count)
	}
}

func TestInitAPP_CascadeDelete(t *testing.T) {
	setupTestDB(t)
	db, _ := InitAPP()

	var userID int
	err := db.QueryRow("INSERT INTO users (username, password_hash) VALUES ('testuser', 'hash') RETURNING id").Scan(&userID)
	if err != nil {
		t.Fatalf("ユーザーが作成できませんでした: %v", err)
	}

	_, err = db.Exec("INSERT INTO user_profiles (user_id, display_name) VALUES ($1, 'test')", userID)
	if err != nil {
		t.Fatalf("プロフィールが作成できませんでした: %v", err)
	}

	_, err = db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		t.Fatalf("ユーザーが削除できませんでした: %v", err)
	}

	var profileCount int
	err = db.QueryRow("SELECT COUNT(*) FROM user_profiles WHERE user_id = $1", userID).Scan(&profileCount)
	if err != nil || profileCount != 0 {
		t.Errorf("CASCADE が削除できませんでした: %v", err)
	}
}
