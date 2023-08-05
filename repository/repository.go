package repository

import (
	"context"

	"tzregion/model"
)

type TodoStorage interface {
	CreateTodo(ctx context.Context, todo *model.Todo) error
}
