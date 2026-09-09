package service

import "testing"

type mockUserRepository struct {
	UserRepository
}

func (r *mockUserRepository) UpdateRole(userID int64, newRole string) error {
	return nil
}

func TestUpdateRoleService_Success(t *testing.T) {
	var repo *mockUserRepository
	userService := NewUserService(repo)

	err := userService.UpdateRole(100, 200, "admin")
	if err != nil {
		t.Errorf("expected no err, got %v", err)
	}
}

func TestUpdateRoleService_WrongRole(t *testing.T) {
	var repo *mockUserRepository
	userService := NewUserService(repo)

	err := userService.UpdateRole(100, 200, "super")
	if err == nil {
		t.Errorf("expected err %q, got nil", "無効なタイプ")
	}
}

func TestUpdateRoleService_Self(t *testing.T) {
	var repo *mockUserRepository
	userService := NewUserService(repo)

	err := userService.UpdateRole(100, 100, "user")
	if err == nil {
		t.Errorf("expected err %q, got nil", "自分の権限を変更することができません")
	}
}
