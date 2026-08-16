package handler

import (
	"net/http"
	"strconv"

	"forum/pkg/httpx"
)

// registerStatsRoutes 注册统计相关路由。
func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats", s.stats)
	mux.HandleFunc("GET /api/stats/boards", s.boardStats)
	mux.HandleFunc("GET /api/stats/hot", s.hotTopics)
}

// stats 全局统计：GET /api/stats
func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.Stats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

// boardStats 按版块统计：GET /api/stats/boards
func (s *Server) boardStats(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.BoardStatsList()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

// hotTopics 热门主题：GET /api/stats/hot?limit=10
func (s *Server) hotTopics(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	result, err := s.svc.HotTopics(limit)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
