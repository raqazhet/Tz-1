package service

import (
	"tzregion/repository"

	"go.uber.org/zap"
)

type Service struct {
	Todo ServiceTodo
}

func NewService(repo *repository.Storage, l *zap.Logger) *Service {
	validateService := NewValidateService()
	todoService := NewAllTodoService(repo, l, validateService)
	return &Service{
		Todo: todoService,
	}
}
