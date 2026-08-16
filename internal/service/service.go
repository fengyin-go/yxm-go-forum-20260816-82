// Package service 实现论坛系统的业务逻辑层。
package service

import (
	"forum/internal/config"
	"forum/internal/store"
	"forum/pkg/logger"
)

// Service 业务逻辑层入口，聚合全部业务能力。
type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

// New 创建业务服务实例。
func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}
