package service

import (
	"testing"

	"forum/internal/model"
)

func TestReplyCountAndStatsConsistency(t *testing.T) {
	s := newTestService()
	_, _, tid := s.seedForum(t)

	if _, err := s.CreateReply(model.Reply{TopicID: tid, Content: "   "}); !model.IsValidationError(err) {
		t.Fatalf("空白回帖内容应被拒绝，得到 %v", err)
	}
	topic, err := s.GetTopic(tid)
	if err != nil {
		t.Fatalf("查询主题失败: %v", err)
	}
	if topic.ReplyCount != 0 {
		t.Fatalf("空白回帖被拒绝后回帖数应为 0，得到 %d", topic.ReplyCount)
	}

	reply, err := s.CreateReply(model.Reply{TopicID: tid, Content: "hello"})
	if err != nil {
		t.Fatalf("创建回帖失败: %v", err)
	}
	replies, err := s.ListReplies(tid)
	if err != nil {
		t.Fatalf("列出回帖失败: %v", err)
	}
	if len(replies) != 1 || replies[0].Content != "hello" {
		t.Fatalf("回帖列表异常: %+v", replies)
	}

	stats, err := s.Stats()
	if err != nil {
		t.Fatalf("统计失败: %v", err)
	}
	if stats.ReplyCount != 1 {
		t.Fatalf("全局回帖数应为 1，得到 %d", stats.ReplyCount)
	}

	if err := s.DeleteReply(reply.ID); err != nil {
		t.Fatalf("删除回帖失败: %v", err)
	}
	topic, _ = s.GetTopic(tid)
	if topic.ReplyCount != 0 {
		t.Fatalf("删除回帖后回帖数应为 0，得到 %d", topic.ReplyCount)
	}
	stats, _ = s.Stats()
	if stats.ReplyCount != 0 {
		t.Fatalf("删除回帖后全局回帖数应为 0，得到 %d", stats.ReplyCount)
	}
}
