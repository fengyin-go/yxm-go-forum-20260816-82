package service

import (
	"errors"
	"testing"
	"time"

	"forum/internal/model"
	"forum/internal/store"
)

func TestUserRegistrationAndStats(t *testing.T) {
	s := newTestService()

	if _, err := s.Register(model.User{Username: "   "}); !model.IsValidationError(err) {
		t.Fatalf("空白用户名应被拒绝，得到 %v", err)
	}

	alice, err := s.Register(model.User{Username: "alice"})
	if err != nil {
		t.Fatalf("注册 alice 失败: %v", err)
	}
	if _, err := s.Register(model.User{Username: "alice"}); !errors.Is(err, store.ErrConflict) {
		t.Fatalf("重复用户名应返回冲突，得到 %v", err)
	}
	bob, err := s.Register(model.User{Username: "bob"})
	if err != nil {
		t.Fatalf("注册 bob 失败: %v", err)
	}

	alice.CreatedAt = time.Unix(1000, 0)
	bob.CreatedAt = time.Unix(2000, 0)
	_ = s.store.UpdateUser(alice)
	_ = s.store.UpdateUser(bob)

	users, err := s.ListUsers()
	if err != nil {
		t.Fatalf("列出用户失败: %v", err)
	}
	if len(users) != 2 || users[0].Username != "bob" {
		t.Fatalf("用户列表应倒序且只有 2 人，得到 %+v", users)
	}

	stats, err := s.Stats()
	if err != nil {
		t.Fatalf("统计失败: %v", err)
	}
	if stats.UserCount != 2 {
		t.Fatalf("用户统计应为 2，得到 %d", stats.UserCount)
	}
}
