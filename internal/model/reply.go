package model

import (
	"strings"
	"time"
)

// Reply 主题下的回帖。
type Reply struct {
	ID        string    `json:"id"`
	TopicID   string    `json:"topic_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 规范化并校验回帖字段。
func (r *Reply) Validate() error {
	r.TopicID = strings.TrimSpace(r.TopicID)
	r.UserID = strings.TrimSpace(r.UserID)
	if r.TopicID == "" {
		return NewValidationError("topic_id", "主题不能为空")
	}
	if r.Content == "" {
		return NewValidationError("content", "回帖内容不能为空")
	}
	return nil
}
