package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	todoapplication "starter/internal/application/todo"
	"starter/internal/domain/todo"
)

// TodoHandler exposes the TodoList use-cases over HTTP (Gin).
type TodoHandler struct {
	svc *todoapplication.Service
}

// NewTodoHandler creates a new handler wired to the given application service.
func NewTodoHandler(svc *todoapplication.Service) *TodoHandler {
	return &TodoHandler{svc: svc}
}

// -----------------------------------------------------------------
// Handlers
// -----------------------------------------------------------------

// Create godoc
// POST /api/v1/todo-lists
// Body: { "title": "..." }
func (h *TodoHandler) Create(c *gin.Context) {
	var req todoapplication.CreateTodoListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		h.handleErr(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto)
}

// GetByID godoc
// GET /api/v1/todo-lists/:id
func (h *TodoHandler) GetByID(c *gin.Context) {
	dto, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, dto)
}

// GetAll godoc
// GET /api/v1/todo-lists
func (h *TodoHandler) GetAll(c *gin.Context) {
	dtos, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		h.handleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, dtos)
}

// AddLine godoc
// POST /api/v1/todo-lists/:id/lines
// Body: { "description": "..." }
func (h *TodoHandler) AddLine(c *gin.Context) {
	var req todoapplication.AddLineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto, err := h.svc.AddLine(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		h.handleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, dto)
}

// CompleteLine godoc
// PATCH /api/v1/todo-lists/:id/lines/:lineId/complete
func (h *TodoHandler) CompleteLine(c *gin.Context) {
	dto, err := h.svc.CompleteLine(c.Request.Context(), c.Param("id"), c.Param("lineId"))
	if err != nil {
		h.handleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, dto)
}

// -----------------------------------------------------------------
// Error mapping
// -----------------------------------------------------------------

func (h *TodoHandler) handleErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, todo.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, todo.ErrLineNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, todo.ErrAlreadyCompleted):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, todo.ErrEmptyTitle),
		errors.Is(err, todo.ErrEmptyDescription):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
