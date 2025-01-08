package handler

import (
	rediscache "github.com/golanguzb70/redis-cache"
	"yalp_ulab/config"
	"yalp_ulab/internal/usecase"
	"yalp_ulab/pkg/logger"
)

type Handler struct {
	Logger  *logger.Logger
	Config  *config.Config
	UseCase *usecase.UseCase
	Redis   rediscache.RedisCache
}

func NewHandler(l *logger.Logger, c *config.Config, useCase *usecase.UseCase, redis rediscache.RedisCache) *Handler {
	return &Handler{
		Logger:  l,
		Config:  c,
		UseCase: useCase,
		Redis:   redis,
	}
}
