package handler

import (
	"context"
	"net/http"
	"time"

	"tzregion/model"
	"tzregion/utils"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Title    string          `json:"title" validate:"required,max=200"`
	ActiveAt utils.CivilTime `json:"activeAt" validate:"required,gte"`
}

// @Summary Create a new todo item
// @Description Create a new todo item with the provided details
// @Tags Todo
// @Accept json
// @Produce json
// @Param todo body Todo true "Todo object"
// @Success 204 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /api/todo-list/tasks [post]
func (h *Handler) CreateTodo(ctx *gin.Context) {
	var todo Todo
	if err := ctx.BindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body", Details: err.Error()})
		return
	}
	ctxx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeOut)
	defer cancel()

	todoModel := model.Todo{
		Title:     todo.Title,
		ActiveAt:  time.Time(todo.ActiveAt),
		CreatedAt: time.Now(),
		Status:    "active",
	}
	if err := h.service.Todo.CreateTodo(ctxx, &todoModel); err != nil {
		ctx.JSON(http.StatusNotFound, model.ErrorResponse{Error: "failed to create", Details: err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, model.SuccessResponse{Message: "todo item created succesfully"})
}

// @Summary Update a todo item
// @Description Update a todo item with the provided ID and details
// @Tags Todo
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param todo body Todo true "Updated todo object"
// @Success 204 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /api/todo-list/tasks/{id} [put]
func (h *Handler) UpdateTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	var todo Todo
	if err := ctx.BindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body", Details: err.Error()})
		return
	}
	ctxx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeOut)
	defer cancel()
	todoModel := model.Todo{
		Title:    todo.Title,
		ActiveAt: time.Time(todo.ActiveAt),
	}
	if err := h.service.Todo.UpdateTodoById(ctxx, id, &todoModel); err != nil {
		ctx.JSON(http.StatusNotFound, model.ErrorResponse{Error: "update todo err", Details: err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, model.SuccessResponse{Message: "todo updated"})
}

// @Summary Delete a todo item by ID
// @Description Delete a todo item with the provided ID
// @Tags Todo
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 204 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/todo-list/tasks/{id} [delete]
func (h *Handler) DeleteTodoById(ctx *gin.Context) {
	id := ctx.Param("id")
	ctxx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeOut)
	defer cancel()
	if err := h.service.Todo.DeleteTodoById(ctxx, id); err != nil {
		ctx.JSON(http.StatusNotFound, model.ErrorResponse{Error: "delete err", Details: err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, model.SuccessResponse{Message: "delete todo by id"})
}

// @Summary Mark a todo item as done
// @Description Mark a todo item with the provided ID as done or not done
// @Tags Todo
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param done path string true "Todo status (done or not done)"
// @Success 204 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/todo-list/tasks/{id}/status/{done} [put]
func (h *Handler) MarkAsDone(ctx *gin.Context) {
	id := ctx.Param("id")
	status := ctx.Param("done")
	ctxx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeOut)
	defer cancel()
	if err := h.service.Todo.MarkAsDone(ctxx, id, status); err != nil {
		ctx.JSON(http.StatusNotFound, model.ErrorResponse{Error: "service err", Details: err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, model.SuccessResponse{Message: "todo succesfully done"})
}

// @Summary Get all todo items by status
// @Description Get a list of all todo items based on the provided status
// @Tags Todo
// @Accept json
// @Produce json
// @Param status path string true "Todo status (active, done, or other)"
// @Success 200 {object} []model.Todo
// @Failure 404 {object} model.ErrorResponse
// @Router /api/todo-list/tasks/{status} [get]
func (h *Handler) FindAllTodos(ctx *gin.Context) {
	status := ctx.Param("status")
	if status == "" {
		status = "active"
	}
	ctxx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeOut)
	defer cancel()
	todos, err := h.service.Todo.FindAll(ctxx, status)
	if err != nil {
		ctx.JSON(http.StatusNotFound, model.ErrorResponse{Error: "not found todos", Details: err.Error()})
		return
	}
	if todos == nil {
		ctx.JSON(http.StatusOK, []*model.Todo{})
	} else {
		ctx.JSON(http.StatusOK, todos)
	}
}
