package repository

import (
	"tzregion/repository/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	mongodb.TodoStorage
}

func NewRepository(mongdd *mongo.Database) *Storage {
	return &Storage{
		TodoStorage: mongodb.NewTodoRepository(mongdd),
	}
}
