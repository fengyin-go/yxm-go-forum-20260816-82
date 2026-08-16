package handler

import (
	"net/http"

	"forum/internal/model"
	"forum/pkg/httpx"
)

// registerBoardRoutes 注册版块相关路由。
func (s *Server) registerBoardRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/boards", s.createBoard)
	mux.HandleFunc("GET /api/boards", s.listBoards)
	mux.HandleFunc("GET /api/boards/{id}", s.getBoard)
	mux.HandleFunc("PUT /api/boards/{id}", s.updateBoard)
	mux.HandleFunc("DELETE /api/boards/{id}", s.deleteBoard)
}

// createBoardRequest 创建版块请求体。
type createBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// createBoard 创建版块：POST /api/boards
func (s *Server) createBoard(w http.ResponseWriter, r *http.Request) {
	var req createBoardRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.CreateBoard(model.Board{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, b)
}

// listBoards 版块列表：GET /api/boards
func (s *Server) listBoards(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.ListBoards()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

// getBoard 版块详情：GET /api/boards/{id}
func (s *Server) getBoard(w http.ResponseWriter, r *http.Request) {
	b, err := s.svc.GetBoard(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

// updateBoardRequest 更新版块请求体。
type updateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// updateBoard 更新版块：PUT /api/boards/{id}
func (s *Server) updateBoard(w http.ResponseWriter, r *http.Request) {
	var req updateBoardRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.UpdateBoard(r.PathValue("id"), model.Board{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

// deleteBoard 删除版块：DELETE /api/boards/{id}
func (s *Server) deleteBoard(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteBoard(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}
