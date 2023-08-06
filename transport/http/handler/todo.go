package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateTodo(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}
