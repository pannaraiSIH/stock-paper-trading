package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db"
)

type HealthHandler struct {
	store *db.Store
}

func NewHealthHandler(store *db.Store) *HealthHandler {
	return &HealthHandler{
		store: store,
	}
}

func (h *HealthHandler) Health(c *gin.Context) {
	ctx := c.Request.Context()

	if err := h.store.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
