package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nielwyn/inventory-system/internal/database"
	apierrors "github.com/nielwyn/inventory-system/pkg/errors"
	"github.com/nielwyn/inventory-system/pkg/response"
)

type HealthHandler struct {
	db *database.Database
}

func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Health(c *gin.Context) {
	response.Success(c, http.StatusOK, "Service is healthy", gin.H{
		"status": "ok",
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	if err := h.db.Health(); err != nil {
		response.Error(c, http.StatusServiceUnavailable, apierrors.CodeServiceUnavail, "Database is not ready")
		return
	}

	response.Success(c, http.StatusOK, "Service is ready", gin.H{
		"status":   "ok",
		"database": "connected",
	})
}
