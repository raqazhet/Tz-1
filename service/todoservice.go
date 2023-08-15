package service

import (
	"context"
	"errors"
	"time"

	"tzregion/model"
	"tzregion/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

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

type ServiceTodo interface {
	CreateTodo(ctx context.Context, todo *model.Todo) error
	UpdateTodoById(ctx context.Context, Id string, todo *model.Todo) error
	DeleteTodoById(ctx context.Context, id string) error
	MarkAsDone(ctx context.Context, Id string, status string) error
	FindAll(ctx context.Context, status string) ([]*model.Todo, error)
	// GetByTitleAndActiveAt(ctx context.Context, title string, activeAt time.Time) (*model.Todo, error)
}

func (s *AllTodService) CreateTodo(ctx context.Context, todo *model.Todo) error {
	if err := s.validate.validateStruct(todo); err != nil {
		s.l.Error("validate error", zap.Error(err))
		return err
	}
	existingTodo, _ := s.repo.FindByTitle(ctx, todo.Title)
	if existingTodo != nil {
		return errors.New("todo already exists")
	}
	return s.repo.CreateTodo(ctx, todo)
}

func (s *AllTodService) UpdateTodoById(ctx context.Context, Id string, todo *model.Todo) error {
	if err := s.validate.validateStruct(todo); err != nil {
		s.l.Error("validate error", zap.Error(err))
		return err
	}
	// Create a new ObjectID from the string ID
	objectID, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		s.l.Error("convert id toObjectID", zap.Error(err))
		return err
	}
	todo.Status = "active"
	return s.repo.UpdateTodoById(ctx, objectID, todo)
}

func (s *AllTodService) DeleteTodoById(ctx context.Context, id string) error {
	// Create a new ObjectID from the string ID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.l.Error("convert id toObjectID", zap.Error(err))
		return err
	}
	return s.repo.DeleteTodoById(ctx, objectID)
}

func (s *AllTodService) MarkAsDone(ctx context.Context, id, status string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.l.Error("convert id todoObjectID", zap.Error(err))
		return err
	}
	todo, err := s.repo.FindTodoById(ctx, objectID)
	if err != nil {
		return err
	}
	if todo == nil {
		return errors.New("todo not found")
	}
	if status == "" {
		status = "done"
	}
	todo.Status = status
	return s.repo.UpdateTodoById(ctx, objectID, todo)
}

func (s *AllTodService) FindAll(ctx context.Context, status string) ([]*model.Todo, error) {
	todos, err := s.repo.TodoStorage.FindAll(ctx, status)
	if err != nil {
		s.l.Error("findAll service err", zap.Error(err))
		return nil, err
	}
	for _, todo := range todos {
		if isWeekend(todo.ActiveAt) {
			todo.Title = "Выходной - " + todo.Title
		}
	}
	return todos, nil
}

func isWeekend(date time.Time) bool {
	return date.Weekday() == time.Saturday || date.Weekday() == time.Sunday
}
