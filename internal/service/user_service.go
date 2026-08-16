package service

import (
	"sort"
	"time"

	"forum/internal/model"
	"forum/pkg/idgen"
)

// Register 注册用户。
func (s *Service) Register(input model.User) (*model.User, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	u := &model.User{
		ID:        idgen.Hex(),
		Username:  input.Username,
		Nickname:  input.Nickname,
		Email:     input.Email,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateUser(u); err != nil {
		return nil, err
	}
	s.log.Infof("注册用户 %s", u.Username)
	return u, nil
}

// GetUser 按 ID 查询用户。
func (s *Service) GetUser(id string) (*model.User, error) {
	return s.store.GetUser(id)
}

// ListUsers 列出全部用户，按注册时间倒序。
func (s *Service) ListUsers() ([]*model.User, error) {
	list := s.store.ListUsers()
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})
	return list, nil
}

// UpdateUser 更新用户资料。
func (s *Service) UpdateUser(id string, input model.User) (*model.User, error) {
	exist, err := s.store.GetUser(id)
	if err != nil {
		return nil, err
	}
	exist.Nickname = input.Nickname
	exist.Email = input.Email
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateUser(exist); err != nil {
		return nil, err
	}
	return exist, nil
}
