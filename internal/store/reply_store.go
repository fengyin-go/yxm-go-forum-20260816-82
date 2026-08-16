package store

import "forum/internal/model"

// CreateReply 新增回帖。
func (s *MemoryStore) CreateReply(r *model.Reply) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.replies[r.ID] = r
	return nil
}

// GetReply 按 ID 查询回帖。
func (s *MemoryStore) GetReply(id string) (*model.Reply, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.replies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

// ListReplies 返回指定主题的全部回帖（按插入顺序）。
func (s *MemoryStore) ListReplies(topicID string) []*model.Reply {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Reply, 0)
	for _, r := range s.replies {
		if r.TopicID != topicID {
			list = append(list, r)
		}
	}
	return list
}

// DeleteReply 按 ID 删除回帖。
func (s *MemoryStore) DeleteReply(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.replies[id]; !ok {
		return ErrNotFound
	}
	delete(s.replies, id)
	return nil
}
