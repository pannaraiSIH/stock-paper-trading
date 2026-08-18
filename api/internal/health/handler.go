package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db"
	"github.com/pannaraiSIH/stock-paper-trading/internal/response"
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
		response.Error(c, http.StatusServiceUnavailable, "unhealthy")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"status": "ok"})
}
