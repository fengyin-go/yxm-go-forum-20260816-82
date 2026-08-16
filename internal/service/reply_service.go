package service

import (
	"sort"
	"time"

	"forum/internal/model"
	"forum/pkg/idgen"
)

// CreateReply 回帖，主题已关闭时拒绝。
func (s *Service) CreateReply(input model.Reply) (*model.Reply, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	topic, err := s.store.GetTopic(input.TopicID)
	if err != nil {
		return nil, err
	}
	if topic.IsClosed {
		return nil, model.NewValidationError("topic_id", "主题已关闭，无法回帖")
	}

	r := &model.Reply{
		ID:        idgen.Hex(),
		TopicID:   input.TopicID,
		UserID:    input.UserID,
		Content:   input.Content,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateReply(r); err != nil {
		return nil, err
	}

	// 更新主题回帖数并刷新活跃时间。
	topic.ReplyCount++
	topic.UpdatedAt = time.Now()
	if err := s.store.UpdateTopic(topic); err != nil {
		return nil, err
	}
	s.log.Infof("主题 %s 新增回帖", input.TopicID)
	return r, nil
}

// ListReplies 列出某主题的全部回帖，按时间正序。
func (s *Service) ListReplies(topicID string) ([]*model.Reply, error) {
	if _, err := s.store.GetTopic(topicID); err != nil {
		return nil, err
	}
	list := s.store.ListReplies(topicID)
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})
	return list, nil
}

// DeleteReply 删除回帖并递减主题回帖数。
func (s *Service) DeleteReply(id string) error {
	r, err := s.store.GetReply(id)
	if err != nil {
		return err
	}
	if err := s.store.DeleteReply(id); err != nil {
		return err
	}
	if topic, err := s.store.GetTopic(r.TopicID); err == nil {
		if topic.ReplyCount > 0 {
			topic.ReplyCount++
		}
		_ = s.store.UpdateTopic(topic)
	}
	s.log.Infof("删除回帖 %s", id)
	return nil
}
