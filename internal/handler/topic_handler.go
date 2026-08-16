package handler

import (
	"net/http"

	"forum/internal/model"
	"forum/pkg/httpx"
)

// registerTopicRoutes 注册主题相关路由。
func (s *Server) registerTopicRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/topics", s.createTopic)
	mux.HandleFunc("GET /api/topics", s.listTopics)
	mux.HandleFunc("GET /api/topics/{id}", s.getTopic)
	mux.HandleFunc("PUT /api/topics/{id}", s.updateTopic)
	mux.HandleFunc("DELETE /api/topics/{id}", s.deleteTopic)
	mux.HandleFunc("PATCH /api/topics/{id}/flag", s.updateTopicFlags)
	mux.HandleFunc("POST /api/topics/{id}/like", s.likeTopic)
	mux.HandleFunc("DELETE /api/topics/{id}/like", s.unlikeTopic)
}

// createTopicRequest 发帖请求体。
type createTopicRequest struct {
	BoardID string `json:"board_id"`
	UserID  string `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// createTopic 发帖：POST /api/topics
func (s *Server) createTopic(w http.ResponseWriter, r *http.Request) {
	var req createTopicRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateTopic(model.Topic{
		BoardID: req.BoardID,
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

// listTopics 主题列表：GET /api/topics?board_id=&keyword=&featured=&page=&size=
func (s *Server) listTopics(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TopicFilter{
		BoardID:      r.URL.Query().Get("board_id"),
		Keyword:      r.URL.Query().Get("keyword"),
		OnlyFeatured: r.URL.Query().Get("featured") == "true",
	}
	items, total, err := s.svc.ListTopics(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

// getTopic 主题详情（累加浏览数）：GET /api/topics/{id}
func (s *Server) getTopic(w http.ResponseWriter, r *http.Request) {
	t, err := s.svc.GetTopic(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

// updateTopicRequest 编辑主题请求体。
type updateTopicRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// updateTopic 编辑主题：PUT /api/topics/{id}
func (s *Server) updateTopic(w http.ResponseWriter, r *http.Request) {
	var req updateTopicRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateTopic(r.PathValue("id"), model.Topic{Title: req.Title, Content: req.Content})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

// deleteTopic 删除主题：DELETE /api/topics/{id}
func (s *Server) deleteTopic(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteTopic(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}

// updateTopicFlagsRequest 主题标记请求体，字段可选。
type updateTopicFlagsRequest struct {
	Pinned   *bool `json:"pinned"`
	Featured *bool `json:"featured"`
	Closed   *bool `json:"closed"`
}

// updateTopicFlags 更新主题标记：PATCH /api/topics/{id}/flag
func (s *Server) updateTopicFlags(w http.ResponseWriter, r *http.Request) {
	var req updateTopicFlagsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateTopicFlags(r.PathValue("id"), req.Pinned, req.Featured, req.Closed)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

// likeTopicRequest 点赞请求体。
type likeTopicRequest struct {
	UserID string `json:"user_id"`
}

// likeTopic 点赞：POST /api/topics/{id}/like
func (s *Server) likeTopic(w http.ResponseWriter, r *http.Request) {
	var req likeTopicRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.LikeTopic(r.PathValue("id"), req.UserID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

// unlikeTopic 取消点赞：DELETE /api/topics/{id}/like?user_id=
func (s *Server) unlikeTopic(w http.ResponseWriter, r *http.Request) {
	t, err := s.svc.UnlikeTopic(r.PathValue("id"), r.URL.Query().Get("user_id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}
