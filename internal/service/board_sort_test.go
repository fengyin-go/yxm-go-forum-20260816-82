package service

import (
	"errors"
	"testing"

	"forum/internal/model"
	"forum/internal/store"
)

func TestBoardSortAndStats(t *testing.T) {
	s := newTestService()

	if _, err := s.CreateBoard(model.Board{Name: "   "}); !model.IsValidationError(err) {
		t.Fatalf("空白名称版块应被拒绝，得到 %v", err)
	}

	a, err := s.CreateBoard(model.Board{Name: "A", SortOrder: 2})
	if err != nil {
		t.Fatalf("创建版块 A 失败: %v", err)
	}
	b, err := s.CreateBoard(model.Board{Name: "B", SortOrder: 1})
	if err != nil {
		t.Fatalf("创建版块 B 失败: %v", err)
	}
	if _, err := s.CreateBoard(model.Board{Name: "A"}); !errors.Is(err, store.ErrConflict) {
		t.Fatalf("重名版块应返回冲突，得到 %v", err)
	}

	boards, err := s.ListBoards()
	if err != nil {
		t.Fatalf("列出版块失败: %v", err)
	}
	if len(boards) != 2 || boards[0].Name != "B" || boards[1].Name != "A" {
		t.Fatalf("版块排序异常: %+v", boards)
	}

	for i := 0; i < 2; i++ {
		if _, err := s.CreateTopic(model.Topic{BoardID: a.ID, Title: "A 帖"}); err != nil {
			t.Fatalf("发帖失败: %v", err)
		}
	}
	if _, err := s.CreateTopic(model.Topic{BoardID: b.ID, Title: "B 帖"}); err != nil {
		t.Fatalf("发帖失败: %v", err)
	}

	boardStats, err := s.BoardStatsList()
	if err != nil {
		t.Fatalf("版块统计失败: %v", err)
	}
	if len(boardStats) != 2 || boardStats[0].BoardID != a.ID || boardStats[0].TopicCount != 2 {
		t.Fatalf("版块统计排序或计数异常: %+v", boardStats)
	}
}
