package http

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) Router() *gin.Engine {
	router := gin.Default()
	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.POST("/api/todo-list/tasks", s.handler.CreateTodo)
	router.PUT("/api/todo-list/tasks/:id", s.handler.UpdateTodo)
	router.DELETE("/api/todo-list/tasks/:id", s.handler.DeleteTodoById)
	router.PUT("/api/todo-list/tasks/:id/done", s.handler.MarkAsDone)
	router.GET("/api/todo-list/tasks/:status", s.handler.FindAllTodos)
	return router
}
