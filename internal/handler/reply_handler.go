package handler

import (
	"net/http"

	"forum/internal/model"
	"forum/pkg/httpx"
)

// registerReplyRoutes 注册回帖相关路由。
func (s *Server) registerReplyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/topics/{id}/replies", s.createReply)
	mux.HandleFunc("GET /api/topics/{id}/replies", s.listReplies)
	mux.HandleFunc("DELETE /api/replies/{id}", s.deleteReply)
}

// createReplyRequest 回帖请求体。
type createReplyRequest struct {
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

// createReply 回帖：POST /api/topics/{id}/replies
func (s *Server) createReply(w http.ResponseWriter, r *http.Request) {
	var req createReplyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	reply, err := s.svc.CreateReply(model.Reply{
		TopicID: r.PathValue("id"),
		UserID:  req.UserID,
		Content: req.Content,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, reply)
}

// listReplies 回帖列表：GET /api/topics/{id}/replies
func (s *Server) listReplies(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.ListReplies(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

// deleteReply 删除回帖：DELETE /api/replies/{id}
func (s *Server) deleteReply(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteReply(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}
