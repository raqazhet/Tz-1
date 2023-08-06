package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"tzregion/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "todos"

type TodoRepo struct {
	collection *mongo.Collection
}

func NewTodoRepository(db *mongo.Database) *TodoRepo {
	todoRepo := &TodoRepo{
		collection: db.Collection(collectionName),
	}
	// Create a unique index for the "title" field
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := todoRepo.collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Printf("failed to create index for collection: %v", err)
		log.Fatal(err)
	}
	return todoRepo
}

func (r *TodoRepo) CreateTodo(ctx context.Context, todo *model.Todo) error {
	_, err := r.collection.InsertOne(ctx, todo)
	if err != nil {
		log.Printf("failed to create err: %v", err)
		return err
	}
	return nil
}

func (r *TodoRepo) UpdateTodoById(ctx context.Context, Id string, todo *model.Todo) error {
	filter := bson.M{"_id": Id}
	updateQuery := bson.M{"$set": bson.M{"title": todo.Title, "activeAt": todo.ActiveAt, "status": todo.Status}}
	res, err := r.collection.UpdateOne(ctx, filter, updateQuery)
	if err != nil {
		log.Printf("Failed to update todo: %v", err)
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("todo with ID %s not found", Id)
	}
	return nil
}

func (r *TodoRepo) DeleteTodoById(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Failed to delete todo: %v", err)
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("todo with ID %s not found", id)
	}
	return nil
}

// func (r *TodoRepo) UpdateStatusDone(ctx context.Context, Id, status string) error {
// 	filter := bson.M{"_id": Id}
// 	update := bson.M{"$set": bson.M{"status": "done"}}
// 	res, err := r.collection.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		log.Printf("failed to status done: %v", err)
// 		return err
// 	}
// 	if res.MatchedCount == 0 {
// 		return err
// 	}
// 	return nil
// }

// FindAll todolists Where status =Active and activeAt<=time.now()
func (r *TodoRepo) FindAll(ctx context.Context, status string) ([]*model.Todo, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}
	options := options.Find()
	options.SetSort(bson.M{"activeAt": 1})
	cur, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		log.Printf("Failed to find todos: %v", err)
		return nil, err
	}
	var todos []*model.Todo
	if err := cur.All(ctx, &todos); err != nil {
		log.Printf("Failed to decode todos: %v", err)
		return nil, err
	}
	return todos, nil
}

func (r *TodoRepo) FindByTitleAndActiveAt(ctx context.Context, title string, activeAt time.Time) (*model.Todo, error) {
	filter := bson.M{"title": title, "activeAt": activeAt}
	var todo model.Todo
	err := r.collection.FindOne(ctx, filter).Decode(&todo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		return nil, err
	}
	return &todo, nil
}
