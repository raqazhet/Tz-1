package service

import (
	"context"
	"errors"
	"time"

	"tzregion/model"
	"tzregion/repository"

	"go.uber.org/zap"
)

type ServiceTodo interface {
	CreateTodo(ctx context.Context, todo *model.Todo) error
	UpdateTodoById(ctx context.Context, Id string, todo *model.Todo) error
	DeleteTodoById(ctx context.Context, id string) error
	MarkAsDone(ctx context.Context, Id string) error
	FindAll(ctx context.Context, status string) ([]*model.Todo, error)
	// GetByTitleAndActiveAt(ctx context.Context, title string, activeAt time.Time) (*model.Todo, error)
}
type AllTodService struct {
	validate *ValidateService
	repo     *repository.Storage
	l        *zap.Logger
}

func NewAllTodoService(repo *repository.Storage, l *zap.Logger, validate *ValidateService) *AllTodService {
	return &AllTodService{
		validate: validate,
		repo:     repo,
		l:        l,
	}
}

func (s *AllTodService) CreateTodo(ctx context.Context, todo *model.Todo) error {
	if err := s.validate.validateStruct(todo); err != nil {
		s.l.Error("validate error", zap.Error(err))
		return err
	}
	existingTodo, err := s.repo.FindByTitleAndActiveAt(ctx, todo.Title, todo.ActiveAt)
	if err != nil {
		return err
	}
	if existingTodo != nil {
		return errors.New("todo already exists")
	}
	return s.repo.CreateTodo(ctx, todo)
}

func (s *AllTodService) UpdateTodoById(ctx context.Context, Id string, todo *model.Todo) error {
	if err := s.validate.validateVariable(Id, "num"); err != nil {
		s.l.Error("validate error", zap.Error(err))
		return err
	}
	if err := s.validate.validateStruct(todo); err != nil {
		s.l.Error("validate error", zap.Error(err))
		return err
	}
	existingTodo, err := s.repo.FindByTitleAndActiveAt(ctx, todo.Title, todo.ActiveAt)
	if err != nil {
		return err
	}
	if existingTodo != nil && existingTodo.ID == Id {
		return errors.New("todo with the same title and activeAt already exists")
	}
	return s.repo.UpdateTodoById(ctx, Id, todo)
}

func (s *AllTodService) DeleteTodoById(ctx context.Context, id string) error {
	if err := s.validate.validateVariable(id, "num"); err != nil {
		s.l.Error("validate err", zap.Error(err))
		return err
	}
	return s.repo.DeleteTodoById(ctx, id)
}

func (s *AllTodService) MarkAsDone(ctx context.Context, id string) error {
	if err := s.validate.validateVariable(id, "num"); err != nil {
		s.l.Error("validate err", zap.Error(err))
		return err
	}
	todo, err := s.repo.FindByTitleAndActiveAt(ctx, id, time.Time{})
	if err != nil {
		return err
	}
	if todo == nil {
		return errors.New("todo not found")
	}
	todo.Status = "done"
	return s.repo.UpdateTodoById(ctx, id, todo)
}

func (s *AllTodService) FindAll(ctx context.Context, status string) ([]*model.Todo, error) {
	if status != "done" {
		status = "active"
	}
	return s.repo.FindAll(ctx, status)
}
