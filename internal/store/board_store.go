package store

import "forum/internal/model"

// CreateBoard 新增版块，ID 或名称重复时返回 ErrConflict。
func (s *MemoryStore) CreateBoard(b *model.Board) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.boards {
		if exist.ID == b.ID || exist.Name == b.Name {
			return ErrConflict
		}
	}
	s.boards[b.ID] = b
	return nil
}

// GetBoard 按 ID 查询版块。
func (s *MemoryStore) GetBoard(id string) (*model.Board, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.boards[id]
	if !ok {
		return nil, ErrNotFound
	}
	return b, nil
}

// GetBoardByName 按名称查询版块。
func (s *MemoryStore) GetBoardByName(name string) (*model.Board, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, b := range s.boards {
		if b.Name == name {
			return b, nil
		}
	}
	return nil, ErrNotFound
}

// ListBoards 返回全部版块。
func (s *MemoryStore) ListBoards() []*model.Board {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Board, 0, len(s.boards))
	for _, b := range s.boards {
		list = append(list, b)
	}
	return list
}

// UpdateBoard 覆盖保存版块。
func (s *MemoryStore) UpdateBoard(b *model.Board) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.boards[b.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.boards {
		if exist.ID != b.ID && exist.Name == b.Name {
			return ErrConflict
		}
	}
	s.boards[b.ID] = b
	return nil
}

// DeleteBoard 按 ID 删除版块。
func (s *MemoryStore) DeleteBoard(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.boards[id]; !ok {
		return ErrNotFound
	}
	delete(s.boards, id)
	return nil
}
