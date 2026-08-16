package service

import (
	"errors"
	"testing"

	"forum/internal/config"
	"forum/internal/model"
	"forum/internal/store"
	"forum/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func (s *Service) seedForum(t *testing.T) (boardID, userID, topicID string) {
	t.Helper()
	u, err := s.Register(model.User{Username: "alice"})
	if err != nil {
		t.Fatalf("注册用户失败: %v", err)
	}
	b, err := s.CreateBoard(model.Board{Name: "技术交流"})
	if err != nil {
		t.Fatalf("创建版块失败: %v", err)
	}
	topic, err := s.CreateTopic(model.Topic{BoardID: b.ID, UserID: u.ID, Title: "首个帖子", Content: "大家好"})
	if err != nil {
		t.Fatalf("发帖失败: %v", err)
	}
	return b.ID, u.ID, topic.ID
}

func TestService_CreateTopic(t *testing.T) {
	s := newTestService()
	// 版块不存在应返回 ErrNotFound
	if _, err := s.CreateTopic(model.Topic{BoardID: "nope", Title: "x"}); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("期望 ErrNotFound，得到 %v", err)
	}

	_, uid, tid := s.seedForum(t)
	topic, err := s.GetTopic(tid)
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	if topic.ViewCount != 1 {
		t.Fatalf("期望浏览数 1，得到 %d", topic.ViewCount)
	}
	if topic.UserID != uid {
		t.Fatalf("作者不匹配")
	}
}

func TestService_ReplyFlow(t *testing.T) {
	s := newTestService()
	_, _, tid := s.seedForum(t)

	if _, err := s.CreateReply(model.Reply{TopicID: tid, Content: "顶一个"}); err != nil {
		t.Fatalf("回帖失败: %v", err)
	}
	topic, _ := s.GetTopic(tid)
	if topic.ReplyCount != 1 {
		t.Fatalf("期望回帖数 1，得到 %d", topic.ReplyCount)
	}

	// 关闭主题后禁止回帖
	closed := true
	if _, err := s.UpdateTopicFlags(tid, nil, nil, &closed); err != nil {
		t.Fatalf("关闭主题失败: %v", err)
	}
	if _, err := s.CreateReply(model.Reply{TopicID: tid, Content: "来晚了"}); !model.IsValidationError(err) {
		t.Fatalf("期望校验错误，得到 %v", err)
	}
}

func TestService_Like(t *testing.T) {
	s := newTestService()
	_, _, tid := s.seedForum(t)

	topic, err := s.LikeTopic(tid, "u1")
	if err != nil {
		t.Fatalf("点赞失败: %v", err)
	}
	if topic.LikeCount != 1 {
		t.Fatalf("期望点赞数 1，得到 %d", topic.LikeCount)
	}

	// 重复点赞
	if _, err := s.LikeTopic(tid, "u1"); !model.IsValidationError(err) {
		t.Fatalf("期望校验错误，得到 %v", err)
	}

	// 取消点赞
	topic, err = s.UnlikeTopic(tid, "u1")
	if err != nil {
		t.Fatalf("取消点赞失败: %v", err)
	}
	if topic.LikeCount != 0 {
		t.Fatalf("期望点赞数 0，得到 %d", topic.LikeCount)
	}
}

func TestService_Stats(t *testing.T) {
	s := newTestService()
	_, _, tid := s.seedForum(t)
	if _, err := s.CreateReply(model.Reply{TopicID: tid, Content: "hi"}); err != nil {
		t.Fatalf("回帖失败: %v", err)
	}

	stats, err := s.Stats()
	if err != nil {
		t.Fatalf("统计失败: %v", err)
	}
	if stats.UserCount != 1 || stats.BoardCount != 1 || stats.TopicCount != 1 {
		t.Fatalf("统计异常: %+v", stats)
	}

	boardStats, _ := s.BoardStatsList()
	if len(boardStats) == 0 {
		t.Fatalf("版块统计为空")
	}
	if boardStats[0].TopicCount != 1 {
		t.Fatalf("版块主题数异常: %+v", boardStats[0])
	}
}
