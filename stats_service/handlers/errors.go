package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"stats/repository"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	if errors.Is(err, repository.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "stats not found for the requested resource",
		})
		return
	}

	slog.Error("internal server error in stats service", slog.Any("err", err))
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "internal server error",
	})
}
