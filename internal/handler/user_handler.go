package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	userapplication "starter/internal/application/user"
	domainuser "starter/internal/domain/user"
)

// UserHandler exposes the User use-cases over HTTP (Gin).
type UserHandler struct {
	svc *userapplication.Service
}

// NewUserHandler creates a new handler wired to the given application service.
func NewUserHandler(svc *userapplication.Service) *UserHandler {
	return &UserHandler{svc: svc}
}

// -----------------------------------------------------------------
// Handlers
// -----------------------------------------------------------------

// GetUsers godoc
// GET /api/v1/users
func (h *UserHandler) GetUsers(c *gin.Context) {
	dtos, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		h.handleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, dtos)
}

// GetUser godoc
// GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	dto, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.handleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, dto)
}

// CreateUser godoc
// POST /api/v1/users
// Body: { "name": "..." }
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req userapplication.CreateUserRequest
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

// -----------------------------------------------------------------
// Error mapping
// -----------------------------------------------------------------

func (h *UserHandler) handleErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domainuser.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domainuser.ErrEmptyName):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
