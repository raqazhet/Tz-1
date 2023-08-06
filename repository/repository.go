package repository

import (
	"context"
	"time"

	"tzregion/model"
	"tzregion/repository/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

type TodoStorage interface {
	CreateTodo(ctx context.Context, todo *model.Todo) error
	UpdateTodoById(ctx context.Context, Id string, todo *model.Todo) error
	DeleteTodoById(ctx context.Context, id string) error
	FindAll(ctx context.Context, status string) ([]*model.Todo, error)
	FindByTitleAndActiveAt(ctx context.Context, title string, activeAt time.Time) (*model.Todo, error)
}
type Storage struct {
	TodoStorage
}

func NewRepository(mongdd *mongo.Database) *Storage {
	return &Storage{
		TodoStorage: mongodb.NewTodoRepository(mongdd),
	}
}
