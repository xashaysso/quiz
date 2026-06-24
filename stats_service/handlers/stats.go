package handlers

import (
	"net/http"
	"stats/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService service.StatsServiceInterface
}

func NewStatsHandler(statsService service.StatsServiceInterface) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

func (h *StatsHandler) GetUserStats(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id format",
		})
		return
	}

	userStats, err := h.statsService.GetUserStats(ctx, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, userStats)
}

func (h *StatsHandler) GetQuizGlobalStats(c *gin.Context) {
	ctx := c.Request.Context()
	quizIDStr := c.Param("quiz_id")
	quizID, err := strconv.ParseInt(quizIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id format",
		})
		return
	}

	quizStats, err := h.statsService.GetQuizGlobalStats(ctx, quizID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, quizStats)
}
