package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"forum/internal/config"
	"forum/internal/model"
	"forum/internal/service"
	"forum/internal/store"
	"forum/pkg/logger"
)

func newHotTestServer() (*Server, *store.MemoryStore) {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	st := store.NewMemoryStore()
	svc := service.New(st, log, cfg)
	return NewServer(svc, log, cfg), st
}

type topicViewEnvelope struct {
	Data struct {
		ViewCount int `json:"view_count"`
	} `json:"data"`
}

type hotTopicsEnvelope struct {
	Data []struct {
		ID        string `json:"id"`
		ViewCount int    `json:"view_count"`
	} `json:"data"`
}

func TestHotTopicsAndViewCount(t *testing.T) {
	srv, st := newHotTestServer()

	board, err := srv.svc.CreateBoard(model.Board{Name: "技术交流"})
	if err != nil {
		t.Fatalf("创建版块失败: %v", err)
	}
	ids := make([]string, 0, 3)
	for _, title := range []string{"A", "B", "C"} {
		topic, err := srv.svc.CreateTopic(model.Topic{BoardID: board.ID, Title: title, Content: "body"})
		if err != nil {
			t.Fatalf("发帖失败: %v", err)
		}
		ids = append(ids, topic.ID)
	}
	for i, id := range ids {
		topic, _ := st.GetTopic(id)
		topic.ViewCount = []int{3, 1, 2}[i]
		_ = st.UpdateTopic(topic)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/topics/"+ids[0], nil)
	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("查看主题返回 %d，body=%s", rec.Code, rec.Body.String())
	}
	var viewed topicViewEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &viewed); err != nil {
		t.Fatalf("解析查看响应失败: %v", err)
	}
	if viewed.Data.ViewCount != 4 {
		t.Fatalf("查看一次后期望浏览数为 4，得到 %d", viewed.Data.ViewCount)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/stats/hot", nil)
	rec = httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("热门主题返回 %d，body=%s", rec.Code, rec.Body.String())
	}
	var hot hotTopicsEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &hot); err != nil {
		t.Fatalf("解析热门主题失败: %v", err)
	}
	if len(hot.Data) != 3 || hot.Data[0].ID != ids[0] || hot.Data[0].ViewCount != 4 {
		t.Fatalf("热门主题排序或数量异常: %+v", hot.Data)
	}
}
