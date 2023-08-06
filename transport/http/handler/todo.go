package handler

import (
	"net/http"
	"time"

	"tzregion/model"

	"github.com/gin-gonic/gin"
)

type todo struct {
	title    string
	activeAt time.Time
}

func (h *Handler) CreateTodo(ctx *gin.Context) {
	var todo todo
	if err := ctx.BindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	todoModel := model.Todo{
		Title:    todo.title,
		ActiveAt: todo.activeAt,
		Status:   "active",
	}
	if err := h.service.Todo.CreateTodo(ctx, &todoModel); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create todo"})
		return
	}
	ctx.Status(http.StatusNoContent)
}
