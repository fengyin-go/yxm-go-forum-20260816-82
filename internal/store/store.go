// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"forum/internal/model"
)

// 数据访问层常见错误。
var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法。
type Store interface {
	// 用户
	CreateUser(u *model.User) error
	GetUser(id string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	ListUsers() []*model.User
	UpdateUser(u *model.User) error
	DeleteUser(id string) error

	// 版块
	CreateBoard(b *model.Board) error
	GetBoard(id string) (*model.Board, error)
	GetBoardByName(name string) (*model.Board, error)
	ListBoards() []*model.Board
	UpdateBoard(b *model.Board) error
	DeleteBoard(id string) error

	// 主题
	CreateTopic(t *model.Topic) error
	GetTopic(id string) (*model.Topic, error)
	ListTopics() []*model.Topic
	UpdateTopic(t *model.Topic) error
	DeleteTopic(id string) error

	// 回帖
	CreateReply(r *model.Reply) error
	GetReply(id string) (*model.Reply, error)
	ListReplies(topicID string) []*model.Reply
	DeleteReply(id string) error

	// 点赞
	AddLike(topicID, userID string) bool
	HasLiked(topicID, userID string) bool
	RemoveLike(topicID, userID string) bool
}
