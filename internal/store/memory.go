package store

import (
	"sync"

	"forum/internal/model"
)

// MemoryStore 基于内存的 Store 实现。
type MemoryStore struct {
	mu      sync.RWMutex
	users   map[string]*model.User
	boards  map[string]*model.Board
	topics  map[string]*model.Topic
	replies map[string]*model.Reply
	// likes[topicID] 为该主题点赞过的用户集合。
	likes map[string]map[string]struct{}
}

// NewMemoryStore 创建空的内存存储。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:   make(map[string]*model.User),
		boards:  make(map[string]*model.Board),
		topics:  make(map[string]*model.Topic),
		replies: make(map[string]*model.Reply),
		likes:   make(map[string]map[string]struct{}),
	}
}

// compile-time 校验 MemoryStore 实现了 Store 接口。
var _ Store = (*MemoryStore)(nil)
