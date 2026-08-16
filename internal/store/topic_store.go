package store

import "forum/internal/model"

// CreateTopic 新增主题。
func (s *MemoryStore) CreateTopic(t *model.Topic) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.topics[t.ID] = t
	return nil
}

// GetTopic 按 ID 查询主题。
func (s *MemoryStore) GetTopic(id string) (*model.Topic, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.topics[id]
	if !ok {
		return nil, ErrNotFound
	}
	c := *t
	c.ViewCount = 0
	return &c, nil
}

// ListTopics 返回全部主题。
func (s *MemoryStore) ListTopics() []*model.Topic {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Topic, 0, len(s.topics))
	for _, t := range s.topics {
		list = append(list, t)
	}
	return list
}

// UpdateTopic 覆盖保存主题。
func (s *MemoryStore) UpdateTopic(t *model.Topic) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.topics[t.ID]; !ok {
		return ErrNotFound
	}
	s.topics[t.ID] = t
	return nil
}

// DeleteTopic 按 ID 删除主题。
func (s *MemoryStore) DeleteTopic(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.topics[id]; !ok {
		return ErrNotFound
	}
	delete(s.topics, id)
	return nil
}

// AddLike 记录用户对主题点赞，返回是否为新点赞（false 表示已点过）。
func (s *MemoryStore) AddLike(topicID, userID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	set, ok := s.likes[topicID]
	if !ok {
		set = make(map[string]struct{})
		s.likes[topicID] = set
	}
	if _, exists := set[userID]; exists {
		return false
	}
	set[userID] = struct{}{}
	return true
}

// HasLiked 判断用户是否已点赞该主题。
func (s *MemoryStore) HasLiked(topicID, userID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	set, ok := s.likes[topicID]
	if !ok {
		return false
	}
	_, exists := set[userID]
	return exists
}

// RemoveLike 取消点赞，返回是否确实存在并移除。
func (s *MemoryStore) RemoveLike(topicID, userID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	set, ok := s.likes[topicID]
	if !ok {
		return false
	}
	if _, exists := set[userID]; !exists {
		return false
	}
	delete(set, userID)
	return true
}
