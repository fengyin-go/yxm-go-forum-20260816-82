package service

import (
	"sort"
	"time"

	"forum/internal/model"
	"forum/pkg/idgen"
)

// CreateBoard 创建版块。
func (s *Service) CreateBoard(input model.Board) (*model.Board, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	b := &model.Board{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Description: input.Description,
		SortOrder:   input.SortOrder,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateBoard(b); err != nil {
		return nil, err
	}
	s.log.Infof("创建版块 %s", b.Name)
	return b, nil
}

// GetBoard 按 ID 查询版块。
func (s *Service) GetBoard(id string) (*model.Board, error) {
	return s.store.GetBoard(id)
}

// ListBoards 列出全部版块，按排序值升序、创建时间倒序。
func (s *Service) ListBoards() ([]*model.Board, error) {
	list := s.store.ListBoards()
	sort.Slice(list, func(i, j int) bool {
		if list[i].SortOrder != list[j].SortOrder {
			return list[i].SortOrder > list[j].SortOrder
		}
		return list[i].CreatedAt.After(list[j].CreatedAt)
	})
	return list, nil
}

// UpdateBoard 更新版块。
func (s *Service) UpdateBoard(id string, input model.Board) (*model.Board, error) {
	exist, err := s.store.GetBoard(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Description = input.Description
	exist.SortOrder = input.SortOrder
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateBoard(exist); err != nil {
		return nil, err
	}
	s.log.Infof("更新版块 %s", exist.Name)
	return exist, nil
}

// DeleteBoard 删除版块。
func (s *Service) DeleteBoard(id string) error {
	if err := s.store.DeleteBoard(id); err != nil {
		return err
	}
	s.log.Infof("删除版块 %s", id)
	return nil
}
