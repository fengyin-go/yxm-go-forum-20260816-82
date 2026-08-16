package model

import (
	"strings"
	"time"
)

// Topic 主题帖。
type Topic struct {
	ID         string    `json:"id"`
	BoardID    string    `json:"board_id"`
	UserID     string    `json:"user_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	ViewCount  int       `json:"view_count"`
	ReplyCount int       `json:"reply_count"`
	LikeCount  int       `json:"like_count"`
	IsPinned   bool      `json:"is_pinned"`
	IsFeatured bool      `json:"is_featured"`
	IsClosed   bool      `json:"is_closed"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Validate 规范化并校验主题字段。
func (t *Topic) Validate() error {
	t.BoardID = strings.TrimSpace(t.BoardID)
	t.UserID = strings.TrimSpace(t.UserID)
	t.Title = strings.TrimSpace(t.Title)
	t.Content = strings.TrimSpace(t.Content)
	if t.BoardID == "" {
		return NewValidationError("board_id", "所属版块不能为空")
	}
	if t.Title == "" {
		return NewValidationError("title", "标题不能为空")
	}
	return nil
}

// TopicFilter 主题列表筛选条件。
type TopicFilter struct {
	BoardID     string
	Keyword     string // 匹配标题
	OnlyFeatured bool
}

// Match 判断主题是否命中筛选条件。
func (f TopicFilter) Match(t *Topic) bool {
	if f.BoardID != "" && t.BoardID != f.BoardID {
		return false
	}
	if f.OnlyFeatured && !t.IsFeatured {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.Title), k) {
			return false
		}
	}
	return true
}
