package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"tzregion/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TodoStorage interface {
	CreateTodo(ctx context.Context, todo *model.Todo) error
	UpdateTodoById(ctx context.Context, Id primitive.ObjectID, todo *model.Todo) error
	DeleteTodoById(ctx context.Context, id primitive.ObjectID) error
	FindAll(ctx context.Context, status string) ([]*model.Todo, error)
	FindByTitle(ctx context.Context, title string) (*model.Todo, error)
	FindTodoById(ctx context.Context, id primitive.ObjectID) (*model.Todo, error)
}

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
		fmt.Println("failed to create todo err ", err)
		return err
	}
	return nil
}

func (r *TodoRepo) UpdateTodoById(ctx context.Context, Id primitive.ObjectID, todo *model.Todo) error {
	filter := bson.M{"_id": Id}
	updateQuery := bson.M{"$set": bson.M{"title": todo.Title, "activeAt": todo.ActiveAt, "status": todo.Status}}
	res, err := r.collection.UpdateOne(ctx, filter, updateQuery)
	if err != nil {
		log.Printf("Failed to update todo: %v", err)
		return err
	}
	if res.MatchedCount == 0 {
		fmt.Println("todo not found by id")
		return fmt.Errorf("todo with ID %s not found", Id)
	}
	return nil
}

func (r *TodoRepo) DeleteTodoById(ctx context.Context, id primitive.ObjectID) error {
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

func (r *TodoRepo) FindByTitle(ctx context.Context, title string) (*model.Todo, error) {
	filter := bson.M{"title": title}
	var todo model.Todo
	err := r.collection.FindOne(ctx, filter).Decode(&todo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No document found")
			return nil, err
		}
		fmt.Println("Failed to execute query:", err)
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepo) FindTodoById(ctx context.Context, id primitive.ObjectID) (*model.Todo, error) {
	var todo model.Todo
	filter := bson.M{"_id": id}
	if err := r.collection.FindOne(ctx, filter).Decode(&todo); err != nil {
		fmt.Printf("find todo by id: %v", err)
		return nil, err
	}
	return &todo, nil
}

// FindAll todolists Where status =Active and activeAt<=time.now()
func (r *TodoRepo) FindAll(ctx context.Context, status string) ([]*model.Todo, error) {
	filter := bson.M{}
	now := time.Now()
	findOptions := options.Find().SetSort(bson.D{{Key: "createdAt", Value: 1}})
	if status == "active" {
		filter = bson.M{
			"status":   status,
			"activeAt": bson.M{"$lte": now},
		}
	} else {
		filter = bson.M{"status": status}
	}
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		fmt.Println("Failed to execute query:", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var todos []*model.Todo
	for cursor.Next(ctx) {
		var todo model.Todo
		err := cursor.Decode(&todo)
		if err != nil {
			fmt.Println("Failed to decode document:", err)
			return nil, err
		}
		todos = append(todos, &todo)
	}

	if err := cursor.Err(); err != nil {
		fmt.Println("Cursor error:", err)
		return nil, err
	}
	return todos, nil
}
