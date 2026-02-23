package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	// imported for swag annotation resolution
	_ "github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// StatsHandler handles stats HTTP endpoints
type StatsHandler struct {
	statsService service.StatsService
	logger       logger.Logger
}

// NewStatsHandler creates a new StatsHandler
func NewStatsHandler(statsService service.StatsService, logger logger.Logger) *StatsHandler {
	return &StatsHandler{statsService: statsService, logger: logger}
}

// GetGlobalStats godoc
// @Summary Get global statistics
// @Tags stats
// @Accept json
// @Produce json
// @Success 200 {object} dto.GlobalStatsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /stats/global [get]
func (h *StatsHandler) GetGlobalStats(c *gin.Context) {
	stats, err := h.statsService.GetGlobalStats(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, stats)
}
