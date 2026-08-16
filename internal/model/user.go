package model

import (
	"strings"
	"time"
)

// User 论坛用户。
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 规范化并校验用户字段。
func (u *User) Validate() error {
	u.Nickname = strings.TrimSpace(u.Nickname)
	u.Email = strings.TrimSpace(u.Email)
	if u.Username == "" {
		return NewValidationError("username", "用户名不能为空")
	}
	if u.Nickname == "" {
		u.Nickname = u.Username
	}
	return nil
}
