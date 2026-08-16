package store

import (
	"testing"
	"time"

	"forum/internal/model"
)

func testUser(name string) *model.User {
	return &model.User{ID: name + "-id", Username: name, Nickname: name, CreatedAt: time.Now()}
}

func TestMemoryStore_User(t *testing.T) {
	s := NewMemoryStore()
	if err := s.CreateUser(testUser("alice")); err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}
	if err := s.CreateUser(testUser("alice")); err != ErrConflict {
		t.Fatalf("期望 ErrConflict，得到 %v", err)
	}
	got, err := s.GetUserByUsername("alice")
	if err != nil || got.ID != "alice-id" {
		t.Fatalf("按用户名查询失败: %v, %v", got, err)
	}
	if len(s.ListUsers()) != 1 {
		t.Fatalf("期望 1 个用户")
	}
	if err := s.DeleteUser("alice-id"); err != nil {
		t.Fatalf("删除失败: %v", err)
	}
}

func TestMemoryStore_Board(t *testing.T) {
	s := NewMemoryStore()
	b := &model.Board{ID: "b1", Name: "技术交流", SortOrder: 1, CreatedAt: time.Now()}
	if err := s.CreateBoard(b); err != nil {
		t.Fatalf("创建版块失败: %v", err)
	}
	if err := s.CreateBoard(&model.Board{ID: "b2", Name: "技术交流"}); err != ErrConflict {
		t.Fatalf("期望 ErrConflict，得到 %v", err)
	}
	got, err := s.GetBoardByName("技术交流")
	if err != nil || got.ID != "b1" {
		t.Fatalf("按名称查询失败: %v, %v", got, err)
	}
}

func TestMemoryStore_TopicAndReply(t *testing.T) {
	s := NewMemoryStore()
	topic := &model.Topic{ID: "t1", BoardID: "b1", Title: "测试帖", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateTopic(topic); err != nil {
		t.Fatalf("创建主题失败: %v", err)
	}
	if len(s.ListTopics()) != 1 {
		t.Fatalf("期望 1 个主题")
	}

	r := &model.Reply{ID: "r1", TopicID: "t1", Content: "沙发", CreatedAt: time.Now()}
	if err := s.CreateReply(r); err != nil {
		t.Fatalf("创建回帖失败: %v", err)
	}
	if got := len(s.ListReplies("t1")); got != 1 {
		t.Fatalf("期望 1 条回帖，得到 %d", got)
	}

	// 点赞去重
	if !s.AddLike("t1", "u1") {
		t.Fatalf("首次点赞应成功")
	}
	if s.AddLike("t1", "u1") {
		t.Fatalf("重复点赞应返回 false")
	}
	if !s.HasLiked("t1", "u1") {
		t.Fatalf("应已点赞")
	}
	if !s.RemoveLike("t1", "u1") {
		t.Fatalf("取消点赞应成功")
	}
	if s.HasLiked("t1", "u1") {
		t.Fatalf("取消后不应再点赞")
	}
}
