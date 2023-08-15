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

// @Summary Create a new todo
// @Description Create a new todo item
// @Accept json
// @Produce json
// @Param todo body Todo true "Todo object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /todos [post]
func (h *Handler) CreateTodo(ctx *gin.Context) {
	var todo Todo
	if err := ctx.BindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "failed to create error", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Todo item created successfully"})
}

// @Summary Update a todo
// @Description Update an existing todo item
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param todo body Todo true "Todo object"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /todos/{id} [put]
func (h *Handler) UpdateTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	var todo Todo
	if err := ctx.BindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}
	ctxx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeOut)
	defer cancel()
	todoModel := model.Todo{
		Title:    todo.Title,
		ActiveAt: time.Time(todo.ActiveAt),
	}
	if err := h.service.Todo.UpdateTodoById(ctxx, id, &todoModel); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"update todo err: ": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{"succes": "todo updated!"})
}

// @Summary Delete a todo by ID
// @Description Delete an existing todo item by ID
// @Param id path string true "Todo ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Router /todos/{id} [delete]
func (h *Handler) DeleteTodoById(ctx *gin.Context) {
	id := ctx.Param("id")
	ctxx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeOut)
	defer cancel()
	if err := h.service.Todo.DeleteTodoById(ctxx, id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{"succes": "deleted todo by id"})
}

// @Summary Mark a todo as done
// @Description Mark a todo as done or undone
// @Param id path string true "Todo ID"
// @Param done path string true "Done status (true/false)"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Router /todos/{id}/done/{done} [put]
func (h *Handler) MarkAsDone(ctx *gin.Context) {
	id := ctx.Param("id")
	status := ctx.Param("done")
	ctxx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeOut)
	defer cancel()
	if err := h.service.Todo.MarkAsDone(ctxx, id, status); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{"message": "todo successfully done"})
}

// @Summary Find all todos
// @Description Find all todos based on status
// @Param status path string false "Todo status (active/done)"
// @Produce json
// @Success 200 {array} model.Todo
// @Failure 404 {object} map[string]interface{}
// @Router /todos/{status} [get]
func (h *Handler) FindAllTodos(ctx *gin.Context) {
	status := ctx.Param("status")
	if status == "" {
		status = "active"
	}
	ctxx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeOut)
	defer cancel()
	todos, err := h.service.Todo.FindAll(ctxx, status)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"not found todos": err.Error()})
		return
	}
	massiv := [1]string{""}
	if todos == nil {
		ctx.JSON(http.StatusOK, massiv)
	} else {
		ctx.JSON(http.StatusOK, todos)
	}
}
