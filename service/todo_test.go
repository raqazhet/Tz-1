package service

import (
	"context"
	"testing"
	"time"

	"tzregion/model"
	"tzregion/repository"
	mock_repo "tzregion/repository/mocks"
	"tzregion/utils"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func Test_Create_Todo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	validate := NewValidateService()
	l := &zap.Logger{}
	todo_repo := mock_repo.NewMockTodoStorage(ctrl)
	repository := &repository.Storage{
		TodoStorage: todo_repo,
	}
	service := NewAllTodoService(repository, l, validate)
	type args struct {
		ctx  context.Context
		todo *model.Todo
	}
	testCases := []struct {
		name          string
		args          args
		buildStubs    func(todo *model.Todo)
		checkResponse func(err error)
	}{
		// Test cases started
		{
			name: "Todo alredy exist",
			args: args{
				todo: &model.Todo{
					Title:     utils.RandomString(10),
					ActiveAt:  time.Now(),
					CreatedAt: time.Now(),
					Status:    "active",
				},
			},
			buildStubs: func(todo *model.Todo) {
				todo_repo.EXPECT().FindByTitle(gomock.Any(), todo.Title).Times(1).Return(&model.Todo{
					ID:        primitive.NewObjectID(),
					Title:     todo.Title,
					ActiveAt:  todo.ActiveAt,
					CreatedAt: todo.CreatedAt,
				}, nil)
				todo_repo.EXPECT().CreateTodo(gomock.Any(), todo).Times(0).Return(nil)
			},
			checkResponse: func(err error) {
				require.Error(t, err)
			},
		},
		{
			name: "Ok",
			args: args{
				todo: &model.Todo{
					Title:     utils.RandomString(20),
					ActiveAt:  time.Now(),
					CreatedAt: time.Now(),
					Status:    "active",
				},
			},
			buildStubs: func(todo *model.Todo) {
				todo_repo.EXPECT().FindByTitle(gomock.Any(), todo.Title).Times(0).Return(&model.Todo{}, mongo.ErrNoDocuments)
				todo_repo.EXPECT().CreateTodo(gomock.Any(), todo).Times(1).Return(nil)
			},
			checkResponse: func(err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStubs(tt.args.todo)
			err := service.CreateTodo(tt.args.ctx, tt.args.todo)
			tt.checkResponse(err)
		})
	}
}
