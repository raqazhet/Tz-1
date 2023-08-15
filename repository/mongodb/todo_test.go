package mongodb

import (
	"context"
	"log"
	"testing"
	"time"
	"tzregion/model"
	"tzregion/utils"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_Create_Todo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	container, err := createMongoContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)
	client, err := mongo_client(ctx, container)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(ctx)
	todo_repo := NewTodoRepository(client.Database("test"))
	todo := &model.Todo{
		Title:     utils.RandomString(10),
		ActiveAt:  time.Now(),
		CreatedAt: time.Now(),
		Status:    "active",
	}
	err = todo_repo.CreateTodo(ctx, todo)
	require.NoError(t, err)
}

func Test_Update_Todo_Byid(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	container, err := createMongoContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)
	client, err := mongo_client(ctx, container)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	todo_repo := NewTodoRepository(client.Database("test"))
	todo := &model.Todo{
		Title:     utils.RandomString(10),
		ActiveAt:  time.Now(),
		CreatedAt: time.Now(),
		Status:    "active",
	}
	err = todo_repo.CreateTodo(ctx, todo)
	require.NoError(t, err)
	one_todo, err := todo_repo.FindByTitle(ctx, todo.Title)
	require.NoError(t, err)
	require.NotNil(t, one_todo)
	require.NoError(t, err)
	object_id, err := primitive.ObjectIDFromHex(one_todo.ID.Hex())
	require.NoError(t, err)
	one_todo.Title = utils.RandomString(10)
	one_todo.Status = "done"
	err = todo_repo.UpdateTodoById(ctx, object_id, one_todo)
	require.NoError(t, err)
}

func Test_Find_By_Title(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	container, err := createMongoContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)
	client, err := mongo_client(ctx, container)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	todo_repo := NewTodoRepository(client.Database("test"))
	todo := &model.Todo{
		Title:     utils.RandomString(10),
		ActiveAt:  time.Now(),
		CreatedAt: time.Now(),
		Status:    "active",
	}
	err = todo_repo.CreateTodo(ctx, todo)
	require.NoError(t, err)
	one_todo, err := todo_repo.FindByTitle(ctx, todo.Title)
	require.NoError(t, err)
	require.NotNil(t, one_todo)
}

func Test_Delete_Todo_By_id(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	container, err := createMongoContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)
	client, err := mongo_client(ctx, container)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	todo_repo := NewTodoRepository(client.Database("test"))
	todo := &model.Todo{
		Title:     utils.RandomString(10),
		ActiveAt:  time.Now(),
		CreatedAt: time.Now(),
		Status:    "active",
	}
	err = todo_repo.CreateTodo(ctx, todo)
	require.NoError(t, err)
	one_todo, err := todo_repo.FindByTitle(ctx, todo.Title)
	require.NoError(t, err)
	require.NotNil(t, one_todo)
	err = todo_repo.DeleteTodoById(ctx, one_todo.ID)
	require.NoError(t, err)
}

func Test_Find_All(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	container, err := createMongoContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)
	client, err := mongo_client(ctx, container)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	todo_repo := NewTodoRepository(client.Database("test"))
	for i := 0; i < 10; i++ {
		todo := &model.Todo{
			Title:     utils.RandomString(10),
			ActiveAt:  time.Now(),
			CreatedAt: time.Now(),
			Status:    "active",
		}
		err = todo_repo.CreateTodo(ctx, todo)
		require.NoError(t, err)
		if i >= 5 {
			todo := &model.Todo{
				Title:     utils.RandomString(10),
				ActiveAt:  time.Now(),
				CreatedAt: time.Now(),
				Status:    "done",
			}
			err = todo_repo.CreateTodo(ctx, todo)
			require.NoError(t, err)
		}
	}
	todos, err := todo_repo.FindAll(ctx, "active")
	require.NoError(t, err)
	require.NotNil(t, todos)
}
