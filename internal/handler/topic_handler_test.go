package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"forum/internal/config"
	"forum/internal/model"
	"forum/internal/service"
	"forum/internal/store"
	"forum/pkg/logger"
)

func newTopicListTestServer() (*Server, *store.MemoryStore) {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	st := store.NewMemoryStore()
	svc := service.New(st, log, cfg)
	return NewServer(svc, log, cfg), st
}

type topicListEnvelope struct {
	Data struct {
		Items []struct {
			Title string `json:"title"`
		} `json:"items"`
		Pagination struct {
			Page  int `json:"page"`
			Size  int `json:"size"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"data"`
}

func getTopicList(t *testing.T, srv *Server, target string) topicListEnvelope {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET %s 返回状态码 %d，body=%s", target, rec.Code, rec.Body.String())
	}
	var out topicListEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("解析 %s 响应失败: %v", target, err)
	}
	return out
}

func TestTopicListSearchAndPagination(t *testing.T) {
	srv, st := newTopicListTestServer()

	board, err := srv.svc.CreateBoard(model.Board{Name: "技术交流"})
	if err != nil {
		t.Fatalf("创建版块失败: %v", err)
	}
	titles := []string{"Hello World", "second topic", "third topic"}
	for i, title := range titles {
		topic, err := srv.svc.CreateTopic(model.Topic{BoardID: board.ID, Title: title, Content: "body"})
		if err != nil {
			t.Fatalf("发帖失败: %v", err)
		}
		topic.UpdatedAt = time.Unix(1000+int64(i), 0)
		_ = st.UpdateTopic(topic)
	}

	// 关键词匹配应忽略大小写，且 featured=false 时不应过滤普通主题。
	search := getTopicList(t, srv, "/api/topics?keyword=hello&page=1&size=20&featured=false")
	if search.Data.Pagination.Total != 1 || len(search.Data.Items) != 1 {
		t.Fatalf("关键词搜索期望 1 条，得到 total=%d items=%d", search.Data.Pagination.Total, len(search.Data.Items))
	}
	if search.Data.Items[0].Title != "Hello World" {
		t.Fatalf("关键词搜索返回错误主题: %s", search.Data.Items[0].Title)
	}

	// 第二页应返回排序后的第二条，而不是跳到第三条。
	page2 := getTopicList(t, srv, "/api/topics?page=2&size=1")
	if len(page2.Data.Items) != 1 || page2.Data.Items[0].Title != "second topic" {
		t.Fatalf("第二页结果异常: %+v", page2.Data.Items)
	}

	// 超过上限的分页大小应被夹紧到上限，而不是变成空结果。
	oversize := getTopicList(t, srv, "/api/topics?page=1&size=200")
	if oversize.Data.Pagination.Size != 100 || len(oversize.Data.Items) != 3 {
		t.Fatalf("超限分页期望返回 3 条且 size=100，得到 size=%d items=%d", oversize.Data.Pagination.Size, len(oversize.Data.Items))
	}
}
