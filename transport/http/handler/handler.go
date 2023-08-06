package handler

import (
	"time"

	"tzregion/service"

	"go.uber.org/zap"
)

const _defaultContextTimeOut = 5 * time.Second

type Handler struct {
	service *service.Service
	l       *zap.Logger
}

func NewHandler(service *service.Service, l *zap.Logger) *Handler {
	return &Handler{
		service: service,
		l:       l,
	}
}
