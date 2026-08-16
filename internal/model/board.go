package model

import (
	"strings"
	"time"
)

// Board 论坛版块。
type Board struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}

// Validate 规范化并校验版块字段。
func (b *Board) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	b.Description = strings.TrimSpace(b.Description)
	if b.Name == "" {
		return NewValidationError("name", "版块名称不能为空")
	}
	if b.SortOrder < 0 {
		return NewValidationError("sort_order", "排序值不能为负数")
	}
	return nil
}
