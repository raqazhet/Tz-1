package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Router() http.Handler {
	router := gin.Default()
	router.POST("/api/todo-list/tasks", s.handler.CreateTodo)
	router.PUT("/api/todo-list/tasks/:id")
	router.DELETE("/api/todo-list/tasks/:id")
	router.PUT("/api/todo-list/tasks/:id/done")
	router.GET("/api/todo-list/tasks")
	return router
}
