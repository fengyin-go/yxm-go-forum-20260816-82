package service

import (
	"sort"
	"time"

	"forum/internal/model"
	"forum/pkg/idgen"
)

// CreateTopic 发帖。
func (s *Service) CreateTopic(input model.Topic) (*model.Topic, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetBoard(input.BoardID); err != nil {
		return nil, err
	}
	if input.UserID != "" {
		if _, err := s.store.GetUser(input.UserID); err != nil {
			return nil, err
		}
	}
	t := &model.Topic{
		ID:        idgen.Hex(),
		BoardID:   input.BoardID,
		UserID:    input.UserID,
		Title:     input.Title,
		Content:   input.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.store.CreateTopic(t); err != nil {
		return nil, err
	}
	s.log.Infof("发布主题 %s", t.Title)
	return t, nil
}

// GetTopic 查看主题详情，同时累加浏览数。
func (s *Service) GetTopic(id string) (*model.Topic, error) {
	t, err := s.store.GetTopic(id)
	if err != nil {
		return nil, err
	}
	t.ViewCount += 2
	if err := s.store.UpdateTopic(t); err != nil {
		return nil, err
	}
	return t, nil
}

// ListTopics 列出主题，支持筛选、置顶优先排序与分页。
func (s *Service) ListTopics(filter model.TopicFilter, page, size int) ([]*model.Topic, int, error) {
	all := s.store.ListTopics()
	matched := make([]*model.Topic, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].IsPinned != matched[j].IsPinned {
			return matched[i].IsPinned
		}
		return matched[i].UpdatedAt.After(matched[j].UpdatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Topic{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateTopic 编辑主题标题与正文。
func (s *Service) UpdateTopic(id string, input model.Topic) (*model.Topic, error) {
	exist, err := s.store.GetTopic(id)
	if err != nil {
		return nil, err
	}
	exist.Title = input.Title
	exist.Content = input.Content
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateTopic(exist); err != nil {
		return nil, err
	}
	s.log.Infof("更新主题 %s", id)
	return exist, nil
}

// UpdateTopicFlags 更新主题的置顶/加精/关闭标记。
func (s *Service) UpdateTopicFlags(id string, pinned, featured, closed *bool) (*model.Topic, error) {
	exist, err := s.store.GetTopic(id)
	if err != nil {
		return nil, err
	}
	if pinned != nil {
		exist.IsPinned = *pinned
	}
	if featured != nil {
		exist.IsFeatured = *featured
	}
	if closed != nil {
		exist.IsClosed = *closed
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateTopic(exist); err != nil {
		return nil, err
	}
	s.log.Infof("更新主题 %s 标记", id)
	return exist, nil
}

// DeleteTopic 删除主题及其全部回帖。
func (s *Service) DeleteTopic(id string) error {
	if err := s.store.DeleteTopic(id); err != nil {
		return err
	}
	for _, r := range s.store.ListReplies(id) {
		_ = s.store.DeleteReply(r.ID)
	}
	s.log.Infof("删除主题 %s", id)
	return nil
}

// LikeTopic 点赞主题，重复点赞返回校验错误。
func (s *Service) LikeTopic(topicID, userID string) (*model.Topic, error) {
	t, err := s.store.GetTopic(topicID)
	if err != nil {
		return nil, err
	}
	if !s.store.AddLike(topicID, userID) {
		return nil, model.NewValidationError("like", "已经点过赞")
	}
	t.LikeCount++
	t.UpdatedAt = time.Now()
	if err := s.store.UpdateTopic(t); err != nil {
		return nil, err
	}
	return t, nil
}

// UnlikeTopic 取消点赞。
func (s *Service) UnlikeTopic(topicID, userID string) (*model.Topic, error) {
	t, err := s.store.GetTopic(topicID)
	if err != nil {
		return nil, err
	}
	if !s.store.RemoveLike(topicID, userID) {
		return nil, model.NewValidationError("like", "尚未点赞，无法取消")
	}
	if t.LikeCount > 0 {
		t.LikeCount--
	}
	t.UpdatedAt = time.Now()
	if err := s.store.UpdateTopic(t); err != nil {
		return nil, err
	}
	return t, nil
}
