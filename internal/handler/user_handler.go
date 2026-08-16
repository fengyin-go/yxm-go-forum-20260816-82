package handler

import (
	"net/http"

	"forum/internal/model"
	"forum/pkg/httpx"
)

// registerUserRoutes 注册用户相关路由。
func (s *Server) registerUserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/users", s.registerUser)
	mux.HandleFunc("GET /api/users", s.listUsers)
	mux.HandleFunc("GET /api/users/{id}", s.getUser)
	mux.HandleFunc("PUT /api/users/{id}", s.updateUser)
}

// registerUserRequest 注册请求体。
type registerUserRequest struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

// registerUser 注册用户：POST /api/users
func (s *Server) registerUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	u, err := s.svc.Register(model.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, u)
}

// listUsers 用户列表：GET /api/users
func (s *Server) listUsers(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.ListUsers()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

// getUser 用户详情：GET /api/users/{id}
func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	u, err := s.svc.GetUser(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, u)
}

// updateUserRequest 更新用户资料请求体。
type updateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

// updateUser 更新用户：PUT /api/users/{id}
func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	var req updateUserRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	u, err := s.svc.UpdateUser(r.PathValue("id"), model.User{
		Nickname: req.Nickname,
		Email:    req.Email,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, u)
}
